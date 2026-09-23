// AUTHORED optional projection binding. Full composition requires migration73;
// narrow profiles retain the unchanged original public model query.
package postgres

func init() {
	electromobilityPublicModelsSQL = `select m.model_id,m.model_code,m.display_name,m.vehicle_class,m.specification from catalog.release_public_models m join platform.tenant t on t.tenant_id=m.tenant_id where t.tenant_code=$1 and t.status='active' order by m.display_name,m.model_id`
}
