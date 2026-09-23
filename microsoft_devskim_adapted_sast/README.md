# Microsoft DevSkim adapted SAST baseline

This pack reconstructs Microsoft DevSkim from the exact signed commit in `source-lock.json`, applies one transparent dependency adaptation, runs the 300 official tests, audits all NuGet graphs, publishes the local CLI and proves one TypeScript finding, one Go finding and one clean fixture.

Run from a newly materialized directory:

```powershell
pwsh -NoProfile -File ./verify_pack.ps1 -AllowNetwork -DotNetExecutable C:/absolute/path/to/dotnet.exe
```

Or build a retained tool and immutable receipt:

```powershell
pwsh -NoProfile -File ./build_and_verify.ps1 -AllowNetwork -DotNetExecutable C:/absolute/path/to/dotnet.exe -Destination C:/absent/devskim-runtime -ReceiptPath C:/absent/devskim-receipt.json
```

The published runtime is `ADAPTED`, not verbatim Microsoft code. It is a no-cost local security-lint baseline for Go and TypeScript. It is not comprehensive interprocedural SAST, does not prove an application secure and does not replace CodeQL entitlement, threat modeling, review, fuzzing, DAST or offensive testing. Re-run the dependency/freshness admission when Microsoft publishes a new release/commit, NuGet advisories change, .NET 10.0.400 leaves support or the project language/claim changes.
