// AUTHORED init adapter. Original worker receives signals in this same process.
using Elite.Cloud.Credentials;
try { CertificateImport.FromEnvironment(); return await Elite.Arca.Wsfe.Worker.Program.Main(args); }
catch { Console.Error.WriteLine("Cloud ARCA initialization failed; controlled diagnostics required."); return 1; }
