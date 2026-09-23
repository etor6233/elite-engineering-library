"""AUTHORED frame adapter; pinned ZXing-C++ only. No image codec or model.
Input: EGY1 + uint16be width + uint16be height + width*height GRAY8 bytes.
Raw decoder output is untrusted. A server session/object owner must authorize it.
"""
import sys,json,hashlib,argparse,struct
from pathlib import Path
MAX_PIXELS=4*1024*1024;MAX_DIMENSION=4096;MAX_TEXT_BYTES=1024;MAX_CODES=8
def emit(status,reason=None,candidates=None,**details):
 print(json.dumps({'schema':'elite-qr-image-decode/v1','status':status,'reason':reason,'candidates':candidates or [],'enabled_symbologies':['QR_CODE','CODE_128','EAN_13'],'authorized':False,'business_effect':False,**details},ensure_ascii=True,separators=(',',':')));return 0
def main():
 p=argparse.ArgumentParser();p.add_argument('--runtime-root',required=True,type=Path);a=p.parse_args();runtime=a.runtime_root.resolve(strict=True)
 raw=sys.stdin.buffer.read(MAX_PIXELS+9)
 if len(raw)<9 or len(raw)>MAX_PIXELS+8 or raw[:4]!=b'EGY1':return emit('REJECTED','FRAME_FORMAT_OR_SIZE')
 width,height=struct.unpack('!HH',raw[4:8])
 if not(20<width<=MAX_DIMENSION and 20<height<=MAX_DIMENSION) or width*height>MAX_PIXELS or len(raw)!=width*height+8:return emit('REJECTED','FRAME_DIMENSIONS')
 sys.path.insert(0,str(runtime));import zxingcpp,importlib.metadata
 if importlib.metadata.version('zxing-cpp')!='3.1.1' or not Path(zxingcpp.__file__).resolve().is_relative_to(runtime):return emit('FAILED','RUNTIME_MISMATCH')
 formats={zxingcpp.BarcodeFormat.QRCode:'QR_CODE',zxingcpp.BarcodeFormat.Code128:'CODE_128',zxingcpp.BarcodeFormat.EAN13:'EAN_13'}
 try:
  raster=memoryview(raw)[8:].cast('B',shape=(height,width))
  decoded=zxingcpp.read_barcodes(raster,formats=zxingcpp.BarcodeFormat.QRCode|zxingcpp.BarcodeFormat.Code128|zxingcpp.BarcodeFormat.EAN13,text_mode=zxingcpp.TextMode.Plain,return_errors=False)
  if len(decoded)>MAX_CODES:return emit('REJECTED','TOO_MANY_CODES')
  candidates=[];seen=set()
  for code in decoded:
   text=code.text
   if not code.valid or not text or len(text.encode('utf8'))>MAX_TEXT_BYTES or code.format not in formats:return emit('REJECTED','INVALID_DECODE_RESULT')
   key=(formats[code.format],text)
   if key in seen:continue
   seen.add(key);points=code.position
   candidates.append({'text':text,'symbology':formats[code.format],'points':[[getattr(points,k).x,getattr(points,k).y] for k in ['top_left','top_right','bottom_right','bottom_left']]})
  status='NOT_FOUND' if not candidates else ('DECODED' if len(candidates)==1 else 'MULTIPLE')
  return emit(status,None,candidates,width=width,height=height,frame_sha256=hashlib.sha256(raw).hexdigest())
 except (ValueError,UnicodeError):return emit('REJECTED','INVALID_DECODE_RESULT')
if __name__=='__main__':raise SystemExit(main())
