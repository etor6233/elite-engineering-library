// AUTHORED Linux container credential glue over .NET X509 APIs; no fiscal logic.
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;

namespace Elite.Cloud.Credentials;
public static class CertificateImport
{
    private static void Require(bool value) { if (!value) throw new InvalidOperationException("Cloud certificate import rejected."); }
    private static string Required(string name) => Environment.GetEnvironmentVariable(name) is { Length: > 0 } value ? value : throw new InvalidOperationException("Cloud certificate configuration missing.");
    private static void CheckFile(string path, long limit)
    {
        Require(Path.IsPathFullyQualified(path));
        for (FileSystemInfo? p = new FileInfo(path); p is not null; p = p is FileInfo f ? f.Directory : ((DirectoryInfo)p).Parent)
            Require(p.LinkTarget is null && !p.Attributes.HasFlag(FileAttributes.ReparsePoint));
        var info = new FileInfo(path); Require(info.Exists && info.Length > 0 && info.Length <= limit);
        if (OperatingSystem.IsLinux()) Require((File.GetUnixFileMode(path) & (UnixFileMode.GroupRead | UnixFileMode.GroupWrite | UnixFileMode.OtherRead | UnixFileMode.OtherWrite)) == 0);
    }
    public static void FromEnvironment()
    {
        if (!OperatingSystem.IsLinux()) throw new InvalidOperationException("Linux container credential bridge required.");
        Require(Required("ARCA_ENVIRONMENT") == "homologation");
        Require(Required("ARCA_CERTIFICATE_STORE") == "CurrentUser");
        var home=Required("HOME"); Require(Path.IsPathFullyQualified(home) && Directory.Exists(home));
        Require((File.GetUnixFileMode(home) & (UnixFileMode.GroupRead | UnixFileMode.GroupWrite | UnixFileMode.OtherRead | UnixFileMode.OtherWrite)) == 0);
        // Cloud Run mounts secret files read-only and may assign root ownership.
        // Copy bounded mounted bytes into this container user's private HOME.
        // Do not require chmod/chown of a provider-managed read-only mount.
        var pfxSource=Required("ARCA_PFX_FILE");var passSource=Required("ARCA_PFX_PASSWORD_FILE");
        Require(pfxSource=="/var/run/secrets/arca/certificate/cert.pfx" && passSource=="/var/run/secrets/arca/password/password");
        var staged=Path.Combine(home,"arca-import-"+Guid.NewGuid());Directory.CreateDirectory(staged);File.SetUnixFileMode(staged,UnixFileMode.UserRead|UnixFileMode.UserWrite|UnixFileMode.UserExecute);
        try
        {
            CopyMounted(pfxSource,Path.Combine(staged,"certificate"),1048576);
            CopyMounted(passSource,Path.Combine(staged,"password"),4096);
            Import(Path.Combine(staged,"certificate"),Path.Combine(staged,"password"),Required("ARCA_CERTIFICATE_THUMBPRINT"));
        }
        finally { Directory.Delete(staged,true); }
    }
    public static void CopyMounted(string source,string destination,int limit)
    {
        if (!OperatingSystem.IsLinux()) throw new InvalidOperationException("Linux secret mount fixture required.");
        var info=new FileInfo(source);Require(info.Exists && info.LinkTarget is null && !info.Attributes.HasFlag(FileAttributes.ReparsePoint) && info.Length>0 && info.Length<=limit);
        using var input=new FileStream(source,FileMode.Open,FileAccess.Read,FileShare.Read);
        using var output=new FileStream(destination,new FileStreamOptions{Mode=FileMode.CreateNew,Access=FileAccess.Write,Share=FileShare.None,UnixCreateMode=UnixFileMode.UserRead|UnixFileMode.UserWrite});
        var bytes=new byte[8192];int count;long total=0;
        try { while((count=input.Read(bytes))>0){total+=count;Require(total<=limit);output.Write(bytes,0,count);}output.Flush(true); }
        finally { CryptographicOperations.ZeroMemory(bytes); }
    }
    public static void Import(string pfx, string passwordFile, string thumbprint)
    {
        CheckFile(pfx, 1048576); CheckFile(passwordFile, 4096);
        Require(thumbprint.Length == 40 && thumbprint.All(Uri.IsHexDigit));
        var password=File.ReadAllText(passwordFile).ToCharArray();
        try
        {
            using var cert=X509CertificateLoader.LoadPkcs12FromFile(pfx,password.AsSpan(),X509KeyStorageFlags.EphemeralKeySet);
            Require(cert.HasPrivateKey && string.Equals(cert.Thumbprint,thumbprint,StringComparison.OrdinalIgnoreCase));
            Require(cert.NotBefore.ToUniversalTime() <= DateTime.UtcNow && cert.NotAfter.ToUniversalTime() > DateTime.UtcNow);
            using var store=new X509Store(StoreName.My, StoreLocation.CurrentUser); store.Open(OpenFlags.ReadWrite);
            store.Add(cert);
            using var installed=store.Certificates.Find(X509FindType.FindByThumbprint,thumbprint,false).Single();
            Require(installed.HasPrivateKey);
            if (OperatingSystem.IsLinux())
            {
                var storePath=Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile),".dotnet/corefx/cryptography/x509stores/my");
                foreach (var path in Directory.EnumerateFiles(storePath))
                    Require((File.GetUnixFileMode(path) & (UnixFileMode.GroupRead | UnixFileMode.GroupWrite | UnixFileMode.OtherRead | UnixFileMode.OtherWrite)) == 0);
            }
        }
        finally { Array.Clear(password); }
    }
}
