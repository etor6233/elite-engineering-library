"""AUTHORED exact evidence contracts for V403 non-production resource adapters.

Names describe checks a target probe/owner must actually perform. This file
does not perform cloud probes or manufacture their results.
"""
from cloud_control import require
GATE_CHECKS={
 "allowed_user_identity":{"authenticated_subject_matches_policy","principal_allowed","tenant_context_bound"},
 "artifact_same_digest":{"registry_digest_matches_receipt","source_manifest_matches","target_platform_matches","family_admissions_valid"},
 "backup_restore":{"restore_completed","restored_schema_matches","restored_data_invariants_pass","recovery_window_met"},
 "budget":{"aggregate_scope_complete","billing_observation_fresh","budget_approval_current","query_cost_limit_declared"},
 "database_user_and_secret_provisioned":{"database_role_least_privilege","secret_version_enabled","connection_target_matches","runtime_user_not_superuser"},
 "dual_authorization_headers":{"workload_token_verified","end_user_token_preserved","issuer_audience_match","cross_organization_access_denied"},
 "empty_database_or_verified_pending_suffix":{"observed_cursor_matches","ordered_migration_hashes_match","unexpected_applied_migrations_absent"},
 "finops_runtime_and_mount_hashes":{"controller_image_digest_matches","gcloud_and_bq_runtime_admitted","runner_source_hash_matches","config_approval_tool_mount_hashes_match","nonroot_private_workdir_verified"},
 "finops_billing_scope_and_query_budget":{"account_table_mapping_verified","declared_project_service_scope_complete","currency_matches","query_byte_limit_approved","export_delay_acknowledged"},
 "finops_seed_generation_zero_and_restore":{"seed_hash_matches","initial_generation_create_only_verified","restored_sqlite_integrity_pass","journal_namespace_matches"},
 "finops_namespace_iam":{"journal_exact_object_scope","intent_create_only_prefix","unrelated_objects_denied","public_access_denied"},
 "finops_schedule_approval_and_reconciliation":{"utc_schedule_matches_interval","approval_current","job_single_task_no_retry","ambiguous_intent_reconciliation_verified","no_automatic_restriction"},
 "finops_pubsub_topic_and_dispatch_approval":{"topic_identity_matches","dispatch_explicitly_approved","subscriber_deduplication_verified","publication_not_human_delivery"},
 "iap_oauth_setup_allowed_users_and_app_oidc_callback":{"iap_enabled","only_allowed_principals","https_origin_preserved","application_oidc_callback_completed","secure_session_cookie_observed"},
 "iap_service_agent_exists":{"service_agent_exists","project_number_matches","invoker_principal_matches"},
 "least_privilege":{"grants_match_declared_policy","unauthorized_operation_denied"},
 "migration_baseline":{"observed_cursor_matches","ordered_migration_hashes_match"},
 "network_access":{"private_route_available","tls_required","public_database_access_denied"},
 "new_resource_absent":{"provider_not_found_observed","project_scope_matches"},
 "no_concurrent_migrator":{"exclusive_migration_lease_held","stale_lease_excluded"},
 "oidc_discovery":{"https_issuer_exact","audience_matches","jwks_refresh_bounded"},
 "pause_resume_reconciliation":{"active_work_drained","webhooks_buffered","scheduler_paused","resume_replay_checked"},
 "private_service_access":{"network_identity_matches","service_allocation_observed","private_sql_access_observed"},
 "retention_policy":{"soft_delete_period_matches","retention_policy_approved","data_class_scope_matches"},
 "retention_scope_and_schedule":{"utc_schedule_approved","tenant_organization_scope_matches","deletion_batch_bounded"},
 "secret_accessor_bindings":{"only_selected_identities_have_access","exact_secret_versions","unauthorized_identity_denied"},
 "secret_versions_provisioned":{"exact_versions_enabled","secret_values_external_to_source"},
 "storage_iam":{"public_access_denied","data_scope_policy_matches","versioning_matches"},
 "version_history_lifecycle_budget":{"versioning_enabled","history_lifecycle_approved","observed_storage_budget_included"},
 "wif_repository_condition":{"issuer_is_github_actions","repository_claim_exact","ref_subject_restricted","impersonation_least_privilege"},
}

def resource_checks(step):
    name=step["id"]
    if step["effect"]=="BARRIER":return {"dependency_receipts_match","all_dependencies_reconciled"}
    if name.startswith("configure-"):return {"resource_identity_matches","image_digests_match","service_account_matches","network_and_mounts_match","scaling_matches","ready_status_observed"}
    if name=="run-migrations":return {"ordered_migration_hashes_match","cursor_reached_selected_end","domain_invariants_pass","no_pending_transaction"}
    if name=="database":return {"sql18_configuration_matches","resource_identity_matches","private_access_only","pitr_backup_policy_matches","ready_status_observed"}
    if name=="database-name":return {"database_identity_matches","owning_instance_matches","database_encoding_matches"}
    if name.startswith("identity-"):return {"principal_identity_matches","project_matches","principal_enabled"}
    if name.startswith(("sql-client-","secret-binding-","web-invoker-")) or name in ("api-invoker","retention-invoker","iap-invoker"):
        return {"binding_resource_matches","principal_and_role_exact","public_grants_absent"}
    if name in {"finops-query-user","finops-billing-reader","finops-journal-access","finops-intent-create","finops-invoker","finops-pubsub-publisher"}:
        return {"binding_resource_matches","principal_and_role_exact","condition_matches_declared_resource_scope","unrelated_resource_access_denied"}
    if name=="finops-journal-role":return {"project_role_identity_matches","permission_set_exact","role_enabled"}
    if name.startswith("storage-versioning-"):return {"bucket_identity_matches","versioning_enabled","history_policy_matches"}
    if name.startswith("storage-"):return {"bucket_identity_matches","retention_and_soft_delete_match","public_access_denied"}
    if name=="iap-api":return {"project_matches","iap_api_enabled"}
    if name in {"scheduler","finops-scheduler"}:return {"job_target_matches","service_account_matches","utc_schedule_matches","retry_limits_match","enabled_state_observed"}
    raise ValueError("No semantic resource contract for step: "+name)

def expected(e,ob):
    matches=[s for s in e["stack"]["commands"] if s["id"]==ob["for_step"]];require(len(matches)==1,"unknown semantic resource")
    step=matches[0]
    if ob["kind"]=="TARGET_GATE_EVIDENCE":
        name=ob["name"];require(name in step["required_target_receipts"] and name in GATE_CHECKS,"unknown or unrelated target gate")
        return step,"TARGET_GATE_VERIFIED",name,GATE_CHECKS[name]
    require(ob["kind"]=="STACK_RESOURCE_RECONCILIATION","semantic observation kind")
    return step,"STACK_RESOURCE_STATE_RECONCILED",step["id"],resource_checks(step)
