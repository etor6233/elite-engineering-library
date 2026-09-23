package whatsappbridge

// AUTHORED qualification glue. Existing owner + real PostgreSQL + Next BFF.
import (
 "encoding/json"
 "os"
 "os/exec"
 "sort"
 "testing"
 "time"
 "elite.local/enterprise/internal/platform/identity"
)

func qualifyFunctionsBrowser(t *testing.T, kind, api, token string, p identity.Principal) {
 t.Helper()
 script:=os.Getenv("ELITE_FUNCTIONS_BROWSER_SCRIPT")
 if script=="" { t.Fatal("explicit browser qualification script required") }
 perms:=[]string{};orgs:=[]string{}
 for k:=range p.Permissions {perms=append(perms,k)}; for k:=range p.Organizations {orgs=append(orgs,k)}
 sort.Strings(perms);sort.Strings(orgs)
 raw,_:=json.Marshal(map[string]any{"subject":p.Subject,"tenantId":p.TenantID,"organizations":orgs,"permissions":perms,"accessToken":token})
 cmd:=exec.Command(os.Getenv("ELITE_WHATSAPP_PYTHON"),script,kind,api,string(raw))
 cmd.Env=os.Environ()
 out,e:=cmd.CombinedOutput();if e!=nil {t.Fatalf("browser failed: %v\n%s",e,out)}
 t.Logf("FUNCTIONS_BROWSER_%s_PASS %s",kind,out)
}

func TestAIActivityNextBrowserConnected(t *testing.T) {
 if os.Getenv("ELITE_FUNCTIONS_BROWSER_SCRIPT")=="" {t.Skip("explicit local browser qualification required")}
 f:=newConnectedReplyFixture(t)
 f.ingest(t,f.inbound(t,"wamid.final-ai-browser",time.Now().Add(-time.Second),"5491112345678"));f.process(t)
 beforeLLM,beforeDomain,beforeMeta:=f.llmCalls.Load(),f.domainCalls.Load(),f.metaCalls.Load()
 qualifyFunctionsBrowser(t,"ai",f.api.URL,"human",f.human)
 if f.llmCalls.Load()!=beforeLLM||f.domainCalls.Load()!=beforeDomain||f.metaCalls.Load()!=beforeMeta {t.Fatal("read UI caused provider effects")}
}
