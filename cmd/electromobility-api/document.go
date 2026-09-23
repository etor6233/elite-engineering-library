package main

// AUTHORED opt-in host composition. FIXTURE uses an in-process no-network HTTP
// transport. PROVIDER requires the original file gate and external AWS credentials.
import (
	"context"
	doc "elite.local/enterprise/internal/documentbridge"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/postgres"
	runtime "example.com/elite/aws-textract-document-runtime/textractruntime"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/jackc/pgx/v5/pgxpool"
)

func init() { documentModuleFactory = selectedDocumentModule }
func selectedDocumentModule(ctx context.Context, pool *pgxpool.Pool, getenv func(string) string) (httpapi.EnterpriseModule, error) {
	if getenv == nil {
		return nil, doc.ErrContract
	}
	enabled := getenv("DOCUMENTS_ENABLED")
	if enabled == "" || enabled == "false" {
		return nil, nil
	}
	if enabled != "true" || pool == nil {
		return nil, doc.ErrContract
	}
	var ready bool
	e := pool.QueryRow(ctx, `select(select count(*)from pg_trigger where not tgisinternal and tgenabled in('O','A')and(tgrelid=to_regclass('document.original')and tgname='original_immutable' or tgrelid=to_regclass('document.extraction')and tgname='extraction_immutable' or tgrelid=to_regclass('document.committed')and tgname in('committed_immutable','document_commit_guard') or tgrelid=to_regclass('approval.request')and tgname='document_review_guard'))=5 and to_regclass('approval.document_review_once')is not null`).Scan(&ready)
	if e != nil || !ready {
		return nil, doc.ErrContract
	}
	mode, profile, hash, root := getenv("DOCUMENTS_MODE"), getenv("DOCUMENTS_PROFILE_FILE"), getenv("DOCUMENTS_PROFILE_SHA256"), getenv("DOCUMENTS_WORK_ROOT")
	var processor doc.Processor
	switch mode {
	case "TYPED_FIXTURE":
		var installed bool
		if err := pool.QueryRow(ctx, `select to_regclass('document.class_schema') is not null and exists(select 1 from pg_trigger where tgrelid=to_regclass('document.original') and tgname='typed_original_guard' and tgenabled in('O','A'))`).Scan(&installed); err != nil || !installed {
			return nil, doc.ErrContract
		}
		processor, e = doc.NewTypedFixturePipeline(profile, hash)
	case "FIXTURE":
		processor, e = doc.NewFixturePipeline(profile, hash, root)
	case "PROVIDER":
		command := doc.SecurityCommand{Python: getenv("DOCUMENTS_PYTHON"), PythonSHA: getenv("DOCUMENTS_PYTHON_SHA256"), Script: getenv("DOCUMENTS_SECURITY_SCRIPT"), ScriptSHA: getenv("DOCUMENTS_SECURITY_SCRIPT_SHA256"), Policy: getenv("DOCUMENTS_SECURITY_POLICY"), PolicySHA: getenv("DOCUMENTS_SECURITY_POLICY_SHA256"), Magika: getenv("DOCUMENTS_MAGIKA"), ClamScan: getenv("DOCUMENTS_CLAMSCAN"), ClamDatabase: getenv("DOCUMENTS_CLAM_DATABASE"), Yara: getenv("DOCUMENTS_YARA"), Rules: getenv("DOCUMENTS_YARA_RULES")}
		var stage doc.SecurityStage
		stage, e = command.Stage()
		if e != nil {
			return nil, e
		}
		region, err := runtime.ApprovedRegion(profile, doc.Class)
		if err != nil {
			return nil, err
		}
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
		if err != nil {
			return nil, doc.ErrContract
		}
		processor, e = doc.NewPipeline(textract.NewFromConfig(cfg, func(o *textract.Options) { o.RetryMaxAttempts = 1 }), stage, profile, hash, root, mode)
	default:
		return nil, doc.ErrContract
	}
	if e != nil {
		return nil, e
	}
	store, e := postgres.NewDocuments(pool, doc.Scope{TenantID: getenv("DOCUMENTS_TENANT_ID"), OrganizationID: getenv("DOCUMENTS_ORGANIZATION_ID"), ProfileSHA: hash, Mode: mode}, processor)
	if e != nil {
		return nil, e
	}
	return &httpapi.DocumentModule{Store: store}, nil
}
