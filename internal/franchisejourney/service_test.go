package franchisejourney

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

type fixedIDs struct{ n int }

type auditedLeadFake struct {
	fakeRepository
	actor string
	calls int
}

func (f *auditedLeadFake) AssignLeadAs(_ context.Context, _, organization, lead, subject string, version int64, _, actor string) (Lead, error) {
	f.actor = actor
	f.calls++
	return Lead{ID: lead, OrganizationID: organization, AssignedSubject: subject, Version: version + 1}, nil
}
func (f *auditedLeadFake) TransitionLeadAs(_ context.Context, _, organization, lead, _, target string, version int64, _, actor string) (Lead, error) {
	f.actor = actor
	f.calls++
	return Lead{ID: lead, OrganizationID: organization, State: target, Version: version + 1}, nil
}

func TestLeadAuditedRepositoryRequired(t *testing.T) {
	ids := &fixedIDs{}
	s := NewService(&fakeRepository{}, ids, fixedClock{})
	if _, err := s.AssignLeadAs(context.Background(), "tenant", "org", "lead", "assignee", 1, "actor"); !errors.Is(err, ErrInvalid) {
		t.Fatal("unadmitted assignment fallback", err)
	}
	if _, err := s.TransitionLeadAs(context.Background(), "tenant", "org", "lead", "new", "contacted", 1, "actor"); !errors.Is(err, ErrInvalid) {
		t.Fatal("unadmitted transition fallback", err)
	}
	if ids.n != 0 {
		t.Fatal("generated IDs before audited admission")
	}
}

func TestLeadActorAndValidation(t *testing.T) {
	repo := &auditedLeadFake{}
	s := NewService(repo, &fixedIDs{}, fixedClock{})
	ctx := context.Background()
	if _, err := s.AssignLeadAs(ctx, "tenant", "org", "lead", "assignee", 1, ""); !errors.Is(err, ErrInvalid) {
		t.Fatal("missing actor", err)
	}
	if _, err := s.TransitionLeadAs(ctx, "tenant", "org", "lead", "new", "contacted", 1, ""); !errors.Is(err, ErrInvalid) {
		t.Fatal("missing transition actor", err)
	}
	if _, err := s.TransitionLeadAs(ctx, "tenant", "org", "lead", "lost", "new", 1, "actor"); !errors.Is(err, ErrInvalid) {
		t.Fatal("illegal transition", err)
	}
	if repo.calls != 0 {
		t.Fatal("invalid command reached persistence")
	}
	v, err := s.AssignLeadAs(ctx, "tenant", "org", "lead", "assignee", 1, "actual-actor")
	if err != nil || v.AssignedSubject != "assignee" || v.Version != 2 || repo.actor != "actual-actor" {
		t.Fatal(v, err, repo.actor)
	}
	v, err = s.TransitionLeadAs(ctx, "tenant", "org", "lead", "new", "contacted", 2, "second-actor")
	if err != nil || v.State != "contacted" || v.Version != 3 || repo.actor != "second-actor" || repo.calls != 2 {
		t.Fatal(v, err, repo.actor)
	}
}

func (i *fixedIDs) New() string { i.n++; return "id-" + string(rune('a'+i.n)) }

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

type fakeRepository struct {
	appointment           Appointment
	slot                  AppointmentSlot
	assigned              Lead
	transition            Lead
	quote                 Quote
	handover              Handover
	accepted              Quote
	resource              ServiceResource
	resourced             Appointment
	appointmentTransition Appointment
	availability          AvailabilityEntry
	checklist             DeliveryChecklist
	completedChecklist    Handover
	deliveryException     DeliveryException
	deliveryResolution    DeliveryResolution
	returnReceipt         ReturnReceipt
	returnDisposition     ReturnDisposition
}

func (f *fakeRepository) PublicLocations(context.Context, string) ([]Location, error) {
	return []Location{{Code: "cordoba"}}, nil
}
func (f *fakeRepository) PublicAppointmentSlots(context.Context, string, string, string, time.Time, time.Time) ([]AppointmentSlot, error) {
	return []AppointmentSlot{{ID: "slot", Capacity: 2, Booked: 1}}, nil
}
func (f *fakeRepository) CreateAppointmentSlot(_ context.Context, _ string, value AppointmentSlot, _ string) (AppointmentSlot, error) {
	f.slot = value
	return value, nil
}
func (f *fakeRepository) RequestAppointment(_ context.Context, _, _, _ string, value Appointment, _, _ string) (Appointment, bool, error) {
	f.appointment = value
	return value, false, nil
}
func (f *fakeRepository) CreateServiceResource(_ context.Context, _ string, value ServiceResource, _ string) (ServiceResource, error) {
	f.resource = value
	return value, nil
}
func (f *fakeRepository) AssignAppointmentResource(_ context.Context, _, _, appointment, resource string, version int64, _ string) (Appointment, error) {
	f.resourced = Appointment{ID: appointment, ResourceID: resource, State: "requested", Version: version + 1}
	return f.resourced, nil
}
func (f *fakeRepository) CreateAvailability(_ context.Context, _, _ string, value AvailabilityEntry, _ string) (AvailabilityEntry, error) {
	f.availability = value
	return value, nil
}
func (f *fakeRepository) CancelAvailability(_ context.Context, _, _, _ string, version int64, _, _, _ string) (AvailabilityEntry, error) {
	return AvailabilityEntry{State: "cancelled", Version: version + 1}, nil
}
func (f *fakeRepository) Availability(context.Context, string, string, string, time.Time, time.Time) ([]AvailabilityEntry, error) {
	return []AvailabilityEntry{{ID: "availability"}}, nil
}
func (f *fakeRepository) TransitionAppointment(_ context.Context, _, _, appointment, _, target string, version int64, _, _, _ string) (Appointment, error) {
	f.appointmentTransition = Appointment{ID: appointment, State: target, Version: version + 1}
	return f.appointmentTransition, nil
}
func (f *fakeRepository) CancelCustomerAppointment(_ context.Context, _, _, _, appointment string, version int64, _, _, _ string) (Appointment, error) {
	return Appointment{ID: appointment, State: "cancelled", Version: version + 1}, nil
}
func (f *fakeRepository) Leads(context.Context, string, string, int, string) (Page[Lead], error) {
	return Page[Lead]{Items: []Lead{{ID: "lead"}}}, nil
}
func (f *fakeRepository) AssignLead(_ context.Context, _, _, _, subject string, version int64, _ string) (Lead, error) {
	f.assigned = Lead{AssignedSubject: subject, Version: version + 1}
	return f.assigned, nil
}
func (f *fakeRepository) TransitionLead(_ context.Context, _, _, _, _, target string, version int64, _ string) (Lead, error) {
	f.transition = Lead{State: target, Version: version + 1}
	return f.transition, nil
}
func (f *fakeRepository) CreateQuote(_ context.Context, _, _ string, value Quote, _, _ string) (Quote, bool, error) {
	f.quote = value
	return value, false, nil
}
func (f *fakeRepository) PublishDeliveryChecklist(_ context.Context, _, _ string, value DeliveryChecklist, _ string) (DeliveryChecklist, error) {
	f.checklist = value
	return value, nil
}
func (f *fakeRepository) CompleteDeliveryChecklist(_ context.Context, _, _, _, _ string, version int64, checklistID string, checklistVersion int64, _ []ChecklistResponse, _ string) (Handover, error) {
	f.completedChecklist = Handover{State: "presented", Version: version + 1, ChecklistID: checklistID, ChecklistVersion: checklistVersion}
	return f.completedChecklist, nil
}
func (f *fakeRepository) RejectHandover(_ context.Context, _, organization, customer, handover string, _ int64, reason, details, _ string, exceptionID, _ string) (DeliveryException, error) {
	f.deliveryException = DeliveryException{ID: exceptionID, OrganizationID: organization, HandoverID: handover, CustomerSubject: customer, ReasonCode: reason, Details: details, State: "open", Version: 1}
	return f.deliveryException, nil
}
func (f *fakeRepository) DeliveryExceptions(context.Context, string, string, int) ([]DeliveryException, error) {
	return []DeliveryException{{ID: "exception", State: "open"}}, nil
}
func (f *fakeRepository) ResolveDeliveryException(_ context.Context, _, _, _ string, exceptionID string, version int64, action, _ string, successorID, authorizationID, _, _ string) (DeliveryResolution, error) {
	f.deliveryResolution = DeliveryResolution{Exception: DeliveryException{ID: exceptionID, State: "resolved", Version: version + 1, ResolutionAction: action}}
	if action == "correct-and-represent" {
		f.deliveryResolution.SuccessorHandover = &Handover{ID: successorID, State: "prepared", Version: 1}
	} else {
		f.deliveryResolution.ReturnAuthorizationID = authorizationID
		f.deliveryResolution.Disposition = action
	}
	return f.deliveryResolution, nil
}
func (f *fakeRepository) ReturnCases(context.Context, string, string, int) ([]ReturnCase, error) {
	return []ReturnCase{{AuthorizationID: "authorization", AuthorizedAction: "return"}}, nil
}
func (f *fakeRepository) ReceiveReturn(_ context.Context, _, organization, subject, authorizationID, serial, condition, notes, evidence, receiptID, _ string) (ReturnReceipt, error) {
	f.returnReceipt = ReturnReceipt{ID: receiptID, AuthorizationID: authorizationID, OrganizationID: organization, ReceivedSerialNumber: serial, ConditionCode: condition, Notes: notes, EvidenceSHA256: evidence, ReceivedBySubject: subject}
	return f.returnReceipt, nil
}
func (f *fakeRepository) DecideReturn(_ context.Context, _, _, subject, receiptID, inventoryAction, notes, dispositionID, inventoryID, remedyID, accountingID, _, _ string) (ReturnDisposition, error) {
	f.returnDisposition = ReturnDisposition{ID: dispositionID, ReceiptID: receiptID, InventoryAction: inventoryAction, CustomerRemedy: "refund", Notes: notes, DecidedBySubject: subject, Effects: []ReturnEffectRequest{{ID: inventoryID, EffectKind: "inventory"}, {ID: remedyID, EffectKind: "refund"}, {ID: accountingID, EffectKind: "accounting"}}}
	return f.returnDisposition, nil
}
func (f *fakeRepository) AcceptQuote(_ context.Context, _, _, _, _ string, version int64, _ string, orderID, _, _, _ string) (Quote, error) {
	f.accepted = Quote{State: "accepted", Version: version + 1, OrderID: orderID}
	return f.accepted, nil
}
func (f *fakeRepository) CustomerJourney(context.Context, string, string, string) (CustomerJourney, error) {
	return CustomerJourney{Quotes: []Quote{{ID: "quote"}}}, nil
}

func (f *fakeRepository) AppointmentAgenda(context.Context, string, string, time.Time, time.Time) (AppointmentAgenda, error) {
	return AppointmentAgenda{Appointments: []Appointment{}, Resources: []ServiceResource{}}, nil
}
func (f *fakeRepository) AcceptHandover(_ context.Context, _, _, _, _ string, version int64, _, checklistID string, checklistVersion int64, evidence, _ string) (Handover, error) {
	f.handover = Handover{State: "accepted", Version: version + 1, AcceptanceEvidence: evidence, ChecklistID: checklistID, ChecklistVersion: checklistVersion}
	return f.handover, nil
}

func newService() (*Service, *fakeRepository, time.Time) {
	now := time.Date(2026, 8, 29, 20, 0, 0, 0, time.UTC)
	repository := &fakeRepository{}
	return NewService(repository, &fixedIDs{}, fixedClock{now}), repository, now
}

func TestPublicAndAppointmentContracts(t *testing.T) {
	service, repository, now := newService()
	if _, err := service.PublicLocations(context.Background(), "Bad_Code"); !errors.Is(err, ErrInvalid) {
		t.Fatal("unsafe tenant code accepted")
	}
	value := Appointment{LeadID: "lead", ModelID: "model", Kind: "test-drive", StartsAt: now.Add(time.Hour)}
	validHash := strings.Repeat("0", 64)
	created, replayed, err := service.RequestAppointment(context.Background(), "tenant", "store", "appointment-key-1", validHash, value)
	if err != nil || replayed || created.State != "requested" || repository.appointment.Version != 1 {
		t.Fatalf("appointment=%+v replayed=%v err=%v", created, replayed, err)
	}
	value.StartsAt = now
	if _, _, err = service.RequestAppointment(context.Background(), "tenant", "store", "appointment-key-2", validHash, value); !errors.Is(err, ErrInvalid) {
		t.Fatal("past appointment accepted")
	}
	value.StartsAt = now.Add(time.Hour)
	if _, _, err = service.RequestAppointment(context.Background(), "tenant", "store", "appointment-key-3", strings.Repeat("z", 64), value); !errors.Is(err, ErrInvalid) {
		t.Fatal("non-hex request digest accepted")
	}
}

func TestAppointmentCapacityContracts(t *testing.T) {
	service, repository, now := newService()
	created, err := service.CreateAppointmentSlot(context.Background(), "tenant", AppointmentSlot{OrganizationID: "store", Kind: "service", StartsAt: now.Add(time.Hour), EndsAt: now.Add(2 * time.Hour), Capacity: 3})
	if err != nil || created.State != "open" || created.Version != 1 || repository.slot.ID == "" {
		t.Fatalf("slot=%+v err=%v", created, err)
	}
	items, err := service.PublicAppointmentSlots(context.Background(), "tenant", "store", "service", now.Add(time.Minute), now.Add(24*time.Hour))
	if err != nil || len(items) != 1 || items[0].Booked != 1 {
		t.Fatalf("slots=%+v err=%v", items, err)
	}
	invalid := []AppointmentSlot{
		{OrganizationID: "store", Kind: "unknown", StartsAt: now.Add(time.Hour), EndsAt: now.Add(2 * time.Hour), Capacity: 1},
		{OrganizationID: "store", Kind: "service", StartsAt: now.Add(time.Hour), EndsAt: now.Add(10 * time.Hour), Capacity: 1},
		{OrganizationID: "store", Kind: "service", StartsAt: now.Add(time.Hour), EndsAt: now.Add(2 * time.Hour), Capacity: 0},
	}
	for _, value := range invalid {
		if _, err = service.CreateAppointmentSlot(context.Background(), "tenant", value); !errors.Is(err, ErrInvalid) {
			t.Fatalf("invalid slot accepted: %+v err=%v", value, err)
		}
	}
	if _, err = service.PublicAppointmentSlots(context.Background(), "tenant", "store", "service", now.Add(time.Minute), now.Add(32*24*time.Hour)); !errors.Is(err, ErrInvalid) {
		t.Fatal("unbounded public slot range accepted")
	}
}

func TestAppointmentResourceAndLifecycleContracts(t *testing.T) {
	service, _, _ := newService()
	resource, err := service.CreateServiceResource(context.Background(), "tenant", ServiceResource{OrganizationID: "store", PrincipalSubject: "technician", DisplayName: "Technician", Kind: "employee", Skills: []string{"service", "delivery"}})
	if err != nil || resource.Status != "active" || resource.Version != 1 || resource.ID == "" {
		t.Fatalf("resource=%+v err=%v", resource, err)
	}
	if _, err = service.CreateServiceResource(context.Background(), "tenant", ServiceResource{OrganizationID: "store", DisplayName: "Invalid employee", Kind: "employee", Skills: []string{"service"}}); !errors.Is(err, ErrInvalid) {
		t.Fatal("employee without principal accepted")
	}
	if _, err = service.CreateServiceResource(context.Background(), "tenant", ServiceResource{OrganizationID: "store", DisplayName: "Bay", Kind: "service-bay", Skills: []string{"service", "service"}}); !errors.Is(err, ErrInvalid) {
		t.Fatal("duplicate resource skill accepted")
	}
	assigned, err := service.AssignAppointmentResource(context.Background(), "tenant", "store", "appointment", resource.ID, 1)
	if err != nil || assigned.ResourceID != resource.ID || assigned.Version != 2 {
		t.Fatalf("assigned=%+v err=%v", assigned, err)
	}
	confirmed, err := service.TransitionAppointment(context.Background(), "tenant", "store", "appointment", "requested", "confirmed", 2, "operator", "")
	if err != nil || confirmed.State != "confirmed" || confirmed.Version != 3 {
		t.Fatalf("confirmed=%+v err=%v", confirmed, err)
	}
	if _, err = service.TransitionAppointment(context.Background(), "tenant", "store", "appointment", "requested", "completed", 2, "operator", ""); !errors.Is(err, ErrInvalid) {
		t.Fatal("appointment skipped confirmation")
	}
}

func TestLeadStateMachineAndAssignment(t *testing.T) {
	service, _, _ := newService()
	assigned, err := service.AssignLead(context.Background(), "tenant", "store", "lead", "sales-person", 1)
	if err != nil || assigned.AssignedSubject != "sales-person" || assigned.Version != 2 {
		t.Fatalf("assignment=%+v err=%v", assigned, err)
	}
	value, err := service.TransitionLead(context.Background(), "tenant", "store", "lead", "new", "contacted", 2)
	if err != nil || value.State != "contacted" || value.Version != 3 {
		t.Fatalf("transition=%+v err=%v", value, err)
	}
	if _, err = service.TransitionLead(context.Background(), "tenant", "store", "lead", "new", "converted", 2); !errors.Is(err, ErrInvalid) {
		t.Fatal("lead skipped states")
	}
}

func TestQuoteAndCustomerOwnershipInputs(t *testing.T) {
	service, _, now := newService()
	quote, replayed, err := service.CreateQuote(context.Background(), "tenant", "quote-request-0001", strings.Repeat("1", 64), Quote{OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", Currency: "ARS", TotalMinorUnits: 1000, ValidUntil: now.Add(24 * time.Hour)})
	if err != nil || replayed || quote.State != "issued" || quote.Version != 1 {
		t.Fatalf("quote=%+v replayed=%v err=%v", quote, replayed, err)
	}
	if _, _, err = service.CreateQuote(context.Background(), "tenant", "quote-request-0002", strings.Repeat("2", 64), Quote{OrganizationID: "store", LeadID: "lead", VariantID: "variant", ValidUntil: now.Add(time.Hour)}); !errors.Is(err, ErrInvalid) {
		t.Fatal("missing price book accepted")
	}
	if _, _, err = service.CreateQuote(context.Background(), "tenant", "short", strings.Repeat("2", 64), Quote{OrganizationID: "store", LeadID: "lead", VariantID: "variant", PriceBookID: "retail", ValidUntil: now.Add(time.Hour)}); !errors.Is(err, ErrInvalid) {
		t.Fatal("short quote idempotency key accepted")
	}
	if _, err = service.CustomerJourney(context.Background(), "tenant", "store", ""); !errors.Is(err, ErrInvalid) {
		t.Fatal("empty customer accepted")
	}
	evidence := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	checklist, err := service.PublishDeliveryChecklist(context.Background(), "tenant", "operator", DeliveryChecklist{ID: "standard-delivery", OrganizationID: "store", Version: 2, Title: "Entrega estándar", Items: []ChecklistItem{{ID: "serial-observed", Prompt: "Verificar serie", ResponseType: "serial", Required: true}, {ID: "asset-condition", Prompt: "Confirmar condición", ResponseType: "confirmation", Required: true}}})
	if err != nil || checklist.State != "published" || len(checklist.Items) != 2 || checklist.Items[1].Ordinal != 2 {
		t.Fatalf("checklist=%+v err=%v", checklist, err)
	}
	completed, err := service.CompleteDeliveryChecklist(context.Background(), "tenant", "store", "operator", "handover", 1, checklist.ID, checklist.Version, []ChecklistResponse{{ItemID: "serial-observed", ResponseText: "SERIAL"}, {ItemID: "asset-condition", ResponseText: "confirmed"}})
	if err != nil || completed.State != "presented" || completed.Version != 2 {
		t.Fatalf("completed=%+v err=%v", completed, err)
	}
	handover, err := service.AcceptHandover(context.Background(), "tenant", "store", "customer", "handover", 2, "SERIAL", checklist.ID, checklist.Version, evidence)
	if err != nil || handover.State != "accepted" {
		t.Fatalf("handover=%+v err=%v", handover, err)
	}
	if _, err = service.AcceptHandover(context.Background(), "tenant", "store", "customer", "handover", 2, "", checklist.ID, checklist.Version, evidence); !errors.Is(err, ErrInvalid) {
		t.Fatalf("handover without serial accepted: %v", err)
	}
	if _, err = service.PublishDeliveryChecklist(context.Background(), "tenant", "operator", DeliveryChecklist{ID: "duplicate", OrganizationID: "store", Version: 1, Title: "Duplicate", Items: []ChecklistItem{{ID: "same", Prompt: "A", ResponseType: "text"}, {ID: "same", Prompt: "B", ResponseType: "text"}}}); !errors.Is(err, ErrInvalid) {
		t.Fatal("duplicate checklist item accepted")
	}
	rejected, err := service.RejectHandover(context.Background(), "tenant", "store", "customer", "handover", 2, "visible-damage", "Rayón en cuadro", evidence)
	if err != nil || rejected.State != "open" || rejected.ReasonCode != "visible-damage" {
		t.Fatalf("rejected=%+v err=%v", rejected, err)
	}
	if _, err = service.RejectHandover(context.Background(), "tenant", "store", "customer", "handover", 2, "Bad_Code", "detail", evidence); !errors.Is(err, ErrInvalid) {
		t.Fatal("invalid rejection reason accepted")
	}
	resolved, err := service.ResolveDeliveryException(context.Background(), "tenant", "store", "operator", rejected.ID, 1, "correct-and-represent", "Corregir preparación")
	if err != nil || resolved.SuccessorHandover == nil || resolved.SuccessorHandover.State != "prepared" {
		t.Fatalf("resolved=%+v err=%v", resolved, err)
	}
	if _, err = service.ResolveDeliveryException(context.Background(), "tenant", "store", "operator", rejected.ID, 1, "refund-now", "invalid"); !errors.Is(err, ErrInvalid) {
		t.Fatal("unknown resolution action accepted")
	}
	receipt, err := service.ReceiveReturn(context.Background(), "tenant", "store", "operator", "authorization", "SERIAL", "damaged", "Recepción con daño", evidence)
	if err != nil || receipt.AuthorizationID != "authorization" || receipt.ConditionCode != "damaged" || receipt.ReceivedBySubject != "operator" {
		t.Fatalf("receipt=%+v err=%v", receipt, err)
	}
	if _, err = service.ReceiveReturn(context.Background(), "tenant", "store", "operator", "authorization", "SERIAL", "unknown", "invalid", evidence); !errors.Is(err, ErrInvalid) {
		t.Fatal("unknown return condition accepted")
	}
	disposition, err := service.DecideReturn(context.Background(), "tenant", "store", "operator", receipt.ID, "quarantine", "Separar hasta efectos downstream")
	if err != nil || disposition.InventoryAction != "quarantine" || disposition.CustomerRemedy != "refund" || len(disposition.Effects) != 3 {
		t.Fatalf("disposition=%+v err=%v", disposition, err)
	}
	if _, err = service.DecideReturn(context.Background(), "tenant", "store", "operator", receipt.ID, "sell-as-new", "invalid"); !errors.Is(err, ErrInvalid) {
		t.Fatal("unknown inventory action accepted")
	}
	accepted, err := service.AcceptQuote(context.Background(), "tenant", "store", "customer", "quote", 1, evidence)
	if err != nil || accepted.State != "accepted" || accepted.OrderID == "" || accepted.Version != 2 {
		t.Fatalf("accepted quote=%+v err=%v", accepted, err)
	}
	if _, err = service.AcceptQuote(context.Background(), "tenant", "store", "customer", "quote", 1, "not-a-digest"); !errors.Is(err, ErrInvalid) {
		t.Fatal("invalid quote acceptance evidence accepted")
	}
}

type checklistReadFake struct {
	fakeRepository
	value DeliveryChecklist
}

func (f *checklistReadFake) PublishedDeliveryChecklist(context.Context, string, string, string, int64) (DeliveryChecklist, error) {
	return f.value, nil
}
func TestPublishedChecklistRejectsIncoherentProjection(t *testing.T) {
	valid := func() DeliveryChecklist {
		return DeliveryChecklist{ID: "checklist", OrganizationID: "store", Version: 1, State: "published", Title: "Synthetic", Items: []ChecklistItem{{ID: "serial", Ordinal: 1, Prompt: "Serie", ResponseType: "serial", Required: true}}}
	}
	cases := map[string]func(*DeliveryChecklist){"id": func(v *DeliveryChecklist) { v.ID = "other" }, "organization": func(v *DeliveryChecklist) { v.OrganizationID = "other" }, "version": func(v *DeliveryChecklist) { v.Version = 2 }, "draft": func(v *DeliveryChecklist) { v.State = "draft" }, "blank title": func(v *DeliveryChecklist) { v.Title = " " }, "long title": func(v *DeliveryChecklist) { v.Title = strings.Repeat("á", 81) }, "no items": func(v *DeliveryChecklist) { v.Items = nil }, "ordinal": func(v *DeliveryChecklist) { v.Items[0].Ordinal = 2 }, "prompt": func(v *DeliveryChecklist) { v.Items[0].Prompt = " " }, "long prompt": func(v *DeliveryChecklist) { v.Items[0].Prompt = strings.Repeat("á", 251) }, "response type": func(v *DeliveryChecklist) { v.Items[0].ResponseType = "unknown" }, "duplicate id": func(v *DeliveryChecklist) { v.Items = append(v.Items, v.Items[0]); v.Items[1].Ordinal = 2 }}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			value := valid()
			mutate(&value)
			s := NewService(&checklistReadFake{value: value}, &fixedIDs{}, fixedClock{})
			got, e := s.PublishedDeliveryChecklist(context.Background(), "tenant", "store", "checklist", 1)
			if !errors.Is(e, ErrConflict) || got.ID != "" {
				t.Fatal(got, e)
			}
		})
	}
	s := NewService(&checklistReadFake{value: valid()}, &fixedIDs{}, fixedClock{})
	if got, e := s.PublishedDeliveryChecklist(context.Background(), "tenant", "store", "checklist", 1); e != nil || got.ID != "checklist" {
		t.Fatal(got, e)
	}
}

type completionReadFake struct {
	fakeRepository
	value ChecklistCompletion
}

func (f *completionReadFake) DeliveryChecklistCompletion(context.Context, string, string, string) (ChecklistCompletion, error) {
	return f.value, nil
}
func TestChecklistCompletionRejectsIncoherentProjection(t *testing.T) {
	valid := func() ChecklistCompletion {
		return ChecklistCompletion{HandoverID: "handover", OrganizationID: "store", State: "presented", Version: 2, ChecklistID: "checklist", ChecklistVersion: 1, CompletedAt: time.Now().UTC(), ActorSubject: "operator", Responses: []ChecklistResponse{{ItemID: "confirmed", ResponseText: "confirmed"}, {ItemID: "serial", ResponseText: "Synthetic"}}}
	}
	cases := map[string]func(*ChecklistCompletion){"handover": func(v *ChecklistCompletion) { v.HandoverID = "other" }, "organization": func(v *ChecklistCompletion) { v.OrganizationID = "other" }, "version": func(v *ChecklistCompletion) { v.Version = 1 }, "prepared": func(v *ChecklistCompletion) { v.State = "prepared" }, "checklist": func(v *ChecklistCompletion) { v.ChecklistID = "" }, "checklist version": func(v *ChecklistCompletion) { v.ChecklistVersion = 0 }, "time": func(v *ChecklistCompletion) { v.CompletedAt = time.Time{} }, "actor": func(v *ChecklistCompletion) { v.ActorSubject = " " }, "no responses": func(v *ChecklistCompletion) { v.Responses = nil }, "order": func(v *ChecklistCompletion) { v.Responses[0], v.Responses[1] = v.Responses[1], v.Responses[0] }, "duplicate": func(v *ChecklistCompletion) { v.Responses[1].ItemID = v.Responses[0].ItemID }, "blank": func(v *ChecklistCompletion) { v.Responses[0].ResponseText = " " }, "too long": func(v *ChecklistCompletion) { v.Responses[0].ResponseText = strings.Repeat("á", 1025) }, "evidence": func(v *ChecklistCompletion) { v.Responses[0].EvidenceSHA256 = "invalid" }}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			v := valid()
			mutate(&v)
			s := NewService(&completionReadFake{value: v}, &fixedIDs{}, fixedClock{})
			got, e := s.DeliveryChecklistCompletion(context.Background(), "tenant", "store", "handover")
			if !errors.Is(e, ErrConflict) || got.HandoverID != "" {
				t.Fatal(got, e)
			}
		})
	}
	for _, state := range []string{"presented", "accepted", "rejected"} {
		v := valid()
		v.State = state
		s := NewService(&completionReadFake{value: v}, &fixedIDs{}, fixedClock{})
		if got, e := s.DeliveryChecklistCompletion(context.Background(), "tenant", "store", "handover"); e != nil || got.State != state {
			t.Fatal(got, e)
		}
	}
}

type returnCaseReadFake struct {
	fakeRepository
	value ReturnCase
}

func (f *returnCaseReadFake) ReturnCaseResult(context.Context, string, string, string) (ReturnCase, error) {
	return f.value, nil
}
func returnResultFixture(action string) ReturnCase {
	now := time.Date(2026, 9, 8, 0, 0, 0, 0, time.UTC)
	remedy := "refund"
	if action == "exchange" {
		remedy = "exchange"
	}
	v := ReturnCase{AuthorizationID: "authorization", OrganizationID: "store", OrderID: "order", StockUnitID: "stock", CustomerSubject: "customer", AuthorizedAction: action, AuthorizedAt: now, Receipt: &ReturnReceipt{ID: "receipt", AuthorizationID: "authorization", OrganizationID: "store", OrderID: "order", StockUnitID: "stock", CustomerSubject: "customer", ReceivedSerialNumber: "SYNTHETIC-SERIAL", ConditionCode: "sealed", Notes: "Synthetic receipt", EvidenceSHA256: strings.Repeat("a", 64), ReceivedBySubject: "receiver", ReceivedAt: now}, Disposition: &ReturnDisposition{ID: "disposition", ReceiptID: "receipt", InventoryAction: "quarantine", CustomerRemedy: remedy, Notes: "Synthetic decision", DecidedBySubject: "decider", DecidedAt: now}}
	requested := []struct{ kind, owner string }{{"inventory", "inventory"}, {remedy, map[string]string{"refund": "payment", "exchange": "fulfillment"}[remedy]}, {"accounting", "accounting"}}
	if action == "return" {
		requested = append(requested, struct{ kind, owner string }{"fiscal", "fiscal"})
	}
	for _, item := range requested {
		id := "request-" + item.kind + "-0001"
		v.Disposition.Effects = append(v.Disposition.Effects, ReturnEffectRequest{ID: id, EffectKind: item.kind, OwnerContext: item.owner, State: "requested", IdempotencyKey: id, RequestedAt: now})
	}
	return v
}
func TestReturnCaseResultRejectsIncoherentProjection(t *testing.T) {
	cases := map[string]func(*ReturnCase){
		"authorization": func(v *ReturnCase) { v.AuthorizationID = "other" }, "organization": func(v *ReturnCase) { v.OrganizationID = "other" }, "action": func(v *ReturnCase) { v.AuthorizedAction = "unknown" }, "authorized time": func(v *ReturnCase) { v.AuthorizedAt = time.Time{} },
		"receipt authorization": func(v *ReturnCase) { v.Receipt.AuthorizationID = "other" }, "receipt organization": func(v *ReturnCase) { v.Receipt.OrganizationID = "other" }, "receipt order": func(v *ReturnCase) { v.Receipt.OrderID = "other" }, "receipt stock": func(v *ReturnCase) { v.Receipt.StockUnitID = "other" }, "receipt customer": func(v *ReturnCase) { v.Receipt.CustomerSubject = "other" }, "serial": func(v *ReturnCase) { v.Receipt.ReceivedSerialNumber = "" }, "condition": func(v *ReturnCase) { v.Receipt.ConditionCode = "unknown" }, "receipt notes": func(v *ReturnCase) { v.Receipt.Notes = " " }, "evidence": func(v *ReturnCase) { v.Receipt.EvidenceSHA256 = "invalid" }, "receiver": func(v *ReturnCase) { v.Receipt.ReceivedBySubject = "" }, "received time": func(v *ReturnCase) { v.Receipt.ReceivedAt = time.Time{} },
		"orphan decision": func(v *ReturnCase) { v.Receipt = nil }, "decision receipt": func(v *ReturnCase) { v.Disposition.ReceiptID = "other" }, "inventory": func(v *ReturnCase) { v.Disposition.InventoryAction = "unknown" }, "remedy": func(v *ReturnCase) { v.Disposition.CustomerRemedy = "exchange" }, "decision notes": func(v *ReturnCase) { v.Disposition.Notes = strings.Repeat("á", 501) }, "decider": func(v *ReturnCase) { v.Disposition.DecidedBySubject = "" }, "decided time": func(v *ReturnCase) { v.Disposition.DecidedAt = time.Time{} },
		"missing effect": func(v *ReturnCase) { v.Disposition.Effects = v.Disposition.Effects[:3] }, "effect owner": func(v *ReturnCase) { v.Disposition.Effects[0].OwnerContext = "payment" }, "effect kind": func(v *ReturnCase) { v.Disposition.Effects[0].EffectKind = "unknown" }, "duplicate id": func(v *ReturnCase) { v.Disposition.Effects[1].ID = v.Disposition.Effects[0].ID }, "duplicate kind": func(v *ReturnCase) { v.Disposition.Effects[1].EffectKind = v.Disposition.Effects[0].EffectKind }, "state": func(v *ReturnCase) { v.Disposition.Effects[0].State = "completed" }, "key": func(v *ReturnCase) { v.Disposition.Effects[0].IdempotencyKey = "short" }, "requested time": func(v *ReturnCase) { v.Disposition.Effects[0].RequestedAt = time.Time{} },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			v := returnResultFixture("return")
			mutate(&v)
			s := NewService(&returnCaseReadFake{value: v}, &fixedIDs{}, fixedClock{})
			got, e := s.ReturnCaseResult(context.Background(), "tenant", "store", "authorization")
			if !errors.Is(e, ErrConflict) || got.AuthorizationID != "" {
				t.Fatal(got, e)
			}
		})
	}
	for _, action := range []string{"return", "exchange"} {
		for _, stage := range []string{"authorized", "received", "decided"} {
			v := returnResultFixture(action)
			if stage == "authorized" {
				v.Receipt = nil
				v.Disposition = nil
			}
			if stage == "received" {
				v.Disposition = nil
			}
			s := NewService(&returnCaseReadFake{value: v}, &fixedIDs{}, fixedClock{})
			if got, e := s.ReturnCaseResult(context.Background(), "tenant", "store", "authorization"); e != nil || got.AuthorizedAction != action {
				t.Fatal(stage, got, e)
			}
		}
	}
}

// Domain seeds model the authorized, received and decided fixtures exercised against PostgreSQL.
// Invariants: reads mint no IDs, never mutate input, expose no partial failure result,
// preserve authorization/receipt scope and retain the existing refund-four/exchange-three request rule.
func FuzzReturnCaseResult(f *testing.F) {
	for _, action := range []string{"return", "exchange"} {
		for _, stage := range []string{"authorized", "received", "decided", "cross-stock"} {
			v := returnResultFixture(action)
			if stage == "authorized" {
				v.Receipt = nil
				v.Disposition = nil
			}
			if stage == "received" {
				v.Disposition = nil
			}
			if stage == "cross-stock" {
				v.Receipt.StockUnitID = "foreign-stock"
			}
			raw, e := json.Marshal(v)
			if e != nil {
				f.Fatal(e)
			}
			f.Add(raw)
		}
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 65536 {
			return
		}
		var v ReturnCase
		if json.Unmarshal(raw, &v) != nil {
			return
		}
		before, e := json.Marshal(v)
		if e != nil {
			return
		}
		ids := &fixedIDs{}
		s := NewService(&returnCaseReadFake{value: v}, ids, fixedClock{})
		got, err := s.ReturnCaseResult(context.Background(), "tenant", "store", "authorization")
		after, e := json.Marshal(v)
		if e != nil || string(before) != string(after) {
			t.Fatal("read mutated its source")
		}
		if ids.n != 0 {
			t.Fatal("read generated mutation IDs")
		}
		if err != nil {
			if !reflect.DeepEqual(got, ReturnCase{}) {
				t.Fatal("error exposed a partial result")
			}
			return
		}
		if got.AuthorizationID != "authorization" || got.OrganizationID != "store" {
			t.Fatal("authorization scope escaped")
		}
		if !reflect.DeepEqual(got, v) {
			t.Fatal("read rewrote source evidence")
		}
		if got.Receipt == nil {
			if got.Disposition != nil {
				t.Fatal("decision without receipt")
			}
			return
		}
		r := got.Receipt
		if r.AuthorizationID != got.AuthorizationID || r.OrganizationID != got.OrganizationID || r.OrderID != got.OrderID || r.StockUnitID != got.StockUnitID || r.CustomerSubject != got.CustomerSubject {
			t.Fatal("receipt binding escaped")
		}
		if got.Disposition == nil {
			return
		}
		d := got.Disposition
		if d.ReceiptID != r.ID {
			t.Fatal("decision binding escaped")
		}
		expected := map[string]string{"inventory": "inventory", "accounting": "accounting"}
		switch got.AuthorizedAction {
		case "return":
			if d.CustomerRemedy != "refund" {
				t.Fatal("refund contract")
			}
			expected["refund"] = "payment"
			expected["fiscal"] = "fiscal"
		case "exchange":
			if d.CustomerRemedy != "exchange" {
				t.Fatal("exchange contract")
			}
			expected["exchange"] = "fulfillment"
		default:
			t.Fatal("unknown authorization action")
		}
		if len(d.Effects) != len(expected) {
			t.Fatal("request count contradicts authorization")
		}
		for _, effect := range d.Effects {
			owner, exists := expected[effect.EffectKind]
			if !exists || effect.OwnerContext != owner || effect.State != "requested" {
				t.Fatal("request contract escaped")
			}
			delete(expected, effect.EffectKind)
		}
		if len(expected) != 0 {
			t.Fatal("missing request owner")
		}
	})
}
