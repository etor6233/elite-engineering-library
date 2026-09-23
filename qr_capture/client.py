"""AUTHORED process boundary. No image codec, URL, path or arbitrary executable
is accepted from capture input. Runtime/Python are trusted deployment settings.
"""
from pathlib import Path
import json,subprocess,sys,threading,struct
_slots=threading.BoundedSemaphore(2);MAX_PIXELS=4*1024*1024
def result(status,reason):return {'schema':'elite-qr-image-decode/v1','status':status,'reason':reason,'candidates':[],'authorized':False,'business_effect':False}
def decode_frame(frame:bytes,runtime_root:Path,python_executable:Path=Path(sys.executable),timeout_seconds:float=5.0):
 if not isinstance(frame,bytes) or not 9<=len(frame)<=MAX_PIXELS+8 or frame[:4]!=b'EGY1':return result('REJECTED','FRAME_FORMAT_OR_SIZE')
 w,h=struct.unpack('!HH',frame[4:8])
 if not(20<w<=4096 and 20<h<=4096) or w*h>MAX_PIXELS or len(frame)!=w*h+8:return result('REJECTED','FRAME_DIMENSIONS')
 if not 0<timeout_seconds<=5:raise ValueError('timeout must be in (0,5] seconds')
 if not _slots.acquire(blocking=False):return result('BUSY','CAPTURE_CAPACITY')
 try:
  args=[str(python_executable.resolve(strict=True)),'-I','-B',str(Path(__file__).with_name('worker.py')),'--runtime-root',str(runtime_root.resolve(strict=True))]
  try:p=subprocess.run(args,input=frame,capture_output=True,timeout=timeout_seconds,check=False)
  except subprocess.TimeoutExpired:return result('TIMEOUT','DECODE_TIMEOUT')
  if p.returncode or len(p.stdout)>65536:return result('FAILED','DECODER_PROCESS')
  try:value=json.loads(p.stdout)
  except (ValueError,UnicodeError):return result('FAILED','DECODER_PROTOCOL')
  if not isinstance(value,dict) or value.get('schema')!='elite-qr-image-decode/v1' or value.get('authorized') is not False or value.get('business_effect') is not False:return result('FAILED','DECODER_PROTOCOL')
  return value
 finally:_slots.release()
