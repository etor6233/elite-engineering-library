// AUTHORED ephemeral self-signed certificate fixture; no user secrets or network.
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;
using Elite.Cloud.Credentials;
if (!OperatingSystem.IsLinux()) throw new InvalidOperationException("Linux fixture only.");
var root=Path.Combine(Path.GetTempPath(),"elite-cloud-cert-"+Guid.NewGuid());Directory.CreateDirectory(root);File.SetUnixFileMode(root,UnixFileMode.UserRead|UnixFileMode.UserWrite|UnixFileMode.UserExecute);
var pfx=Path.Combine(root,"fixture.pfx");var password=Path.Combine(root,"password");
using var key=RSA.Create(2048);var request=new CertificateRequest("CN=V403 synthetic fixture",key,HashAlgorithmName.SHA256,RSASignaturePadding.Pkcs1);
using var cert=request.CreateSelfSigned(DateTimeOffset.UtcNow.AddMinutes(-1),DateTimeOffset.UtcNow.AddMinutes(10));
var pass=Convert.ToHexString(RandomNumberGenerator.GetBytes(32));File.WriteAllBytes(pfx,cert.Export(X509ContentType.Pfx,pass));File.WriteAllText(password,pass);
foreach(var f in new[]{pfx,password})File.SetUnixFileMode(f,UnixFileMode.UserRead|UnixFileMode.UserWrite);
var checks=0;
try
{
    CertificateImport.Import(pfx,password,cert.Thumbprint);checks++;
    var mountedCopy=Path.Combine(root,"copy");File.SetUnixFileMode(pfx,UnixFileMode.UserRead|UnixFileMode.GroupRead|UnixFileMode.OtherRead);
    CertificateImport.CopyMounted(pfx,mountedCopy,1048576);CertificateImport.Import(mountedCopy,password,cert.Thumbprint);checks++;
    File.SetUnixFileMode(pfx,UnixFileMode.UserRead|UnixFileMode.UserWrite);
    using var store=new X509Store(StoreName.My,StoreLocation.CurrentUser);store.Open(OpenFlags.ReadWrite);
    using var loaded=store.Certificates.Find(X509FindType.FindByThumbprint,cert.Thumbprint,false).Single();
    using var privateKey=loaded.GetRSAPrivateKey();var data=RandomNumberGenerator.GetBytes(24);var sig=privateKey!.SignData(data,HashAlgorithmName.SHA256,RSASignaturePadding.Pkcs1);
    if(!key.VerifyData(data,sig,HashAlgorithmName.SHA256,RSASignaturePadding.Pkcs1))throw new Exception("sign fixture failed");checks++;
    void Reject(Action action){try{action();}catch(InvalidOperationException){checks++;return;}catch(CryptographicException){checks++;return;}throw new Exception("negative fixture accepted");}
    Reject(()=>CertificateImport.Import(pfx,password,new string('0',40)));
    File.WriteAllText(password,"wrong fixture value");Reject(()=>CertificateImport.Import(pfx,password,cert.Thumbprint));File.WriteAllText(password,pass);
    File.SetUnixFileMode(password,UnixFileMode.UserRead|UnixFileMode.OtherRead);Reject(()=>CertificateImport.Import(pfx,password,cert.Thumbprint));File.SetUnixFileMode(password,UnixFileMode.UserRead|UnixFileMode.UserWrite);
    var link=Path.Combine(root,"linked");File.CreateSymbolicLink(link,pfx);Reject(()=>CertificateImport.Import(link,password,cert.Thumbprint));
    store.Remove(loaded);
    Console.WriteLine(System.Text.Json.JsonSerializer.Serialize(new{result="PASS",checks,method="EXECUTED_LINUX_SYNTHETIC_CERTIFICATE",cloud_executed=false,live_arca=false}));
}
finally { Directory.Delete(root,true); }
