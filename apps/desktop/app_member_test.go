package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	dbsqlite "github.com/zbango/gym-saas/apps/desktop/internal/sqlite"
	"github.com/zbango/gym-saas/apps/desktop/internal/updater"
	"github.com/zbango/gym-saas/go/core/application"
	"github.com/zbango/gym-saas/go/core/domain"
)

func TestAppMemberMethodsUseTheMemberService(t *testing.T) {
	store, err := dbsqlite.Open(filepath.Join(t.TempDir(), "gym-saas.db"))
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	defer store.Close()
	gymID, err := domain.NewGymID()
	if err != nil {
		t.Fatalf("NewGymID returned error: %v", err)
	}
	now := time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC)
	gym, err := domain.NewGym(gymID, "Test gym", "UTC", now, now)
	if err != nil {
		t.Fatalf("NewGym returned error: %v", err)
	}
	if err := store.EnsureGym(context.Background(), gym); err != nil {
		t.Fatalf("EnsureGym returned error: %v", err)
	}
	service, err := application.NewMemberService(dbsqlite.NewMemberRepository(store), gymID, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewMemberService returned error: %v", err)
	}
	app := NewApp(updater.NewClient("test"), "test", service)
	created, err := app.CreateMember(MemberInput{
		FirstName: "Ada", LastName: "Lovelace", Phone: "555-0100", Email: "ada@example.com",
		DateOfBirth: "1815-12-10", Address: "12 St. James's Square", Status: "active",
	})
	if err != nil {
		t.Fatalf("CreateMember returned error: %v", err)
	}
	if created.ID == "" || created.FirstName != "Ada" {
		t.Fatalf("CreateMember = %#v", created)
	}
	updated, err := app.UpdateMember(created.ID, MemberInput{
		FirstName: "Augusta", LastName: "Lovelace", Phone: "555-0100", Email: "ada@example.com",
		DateOfBirth: "1815-12-10", Address: "12 St. James's Square", Status: "inactive",
	})
	if err != nil {
		t.Fatalf("UpdateMember returned error: %v", err)
	}
	if updated.FirstName != "Augusta" || updated.Status != "inactive" {
		t.Fatalf("UpdateMember = %#v", updated)
	}
	if err := app.ArchiveMember(created.ID); err != nil {
		t.Fatalf("ArchiveMember returned error: %v", err)
	}
	members, err := app.ListMembers()
	if err != nil {
		t.Fatalf("ListMembers returned error: %v", err)
	}
	if len(members) != 0 {
		t.Fatalf("ListMembers length = %d, want 0 after archive", len(members))
	}
}

func TestAppListsMembersAfterStoreReopens(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gym-saas.db")
	store, err := dbsqlite.Open(path)
	if err != nil {
		t.Fatalf("Open returned error: %v", err)
	}
	gymID, err := domain.NewGymID()
	if err != nil {
		t.Fatalf("NewGymID returned error: %v", err)
	}
	now := time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC)
	gym, err := domain.NewGym(gymID, "Test gym", "UTC", now, now)
	if err != nil {
		t.Fatalf("NewGym returned error: %v", err)
	}
	if err := store.EnsureGym(context.Background(), gym); err != nil {
		t.Fatalf("EnsureGym returned error: %v", err)
	}
	service, err := application.NewMemberService(dbsqlite.NewMemberRepository(store), gymID, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewMemberService returned error: %v", err)
	}
	app := NewApp(updater.NewClient("test"), "test", service)
	created, err := app.CreateMember(MemberInput{FirstName: "Ada", LastName: "Lovelace", Phone: "555-0100", Status: "active"})
	if err != nil {
		t.Fatalf("CreateMember returned error: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	store, err = dbsqlite.Open(path)
	if err != nil {
		t.Fatalf("reopen returned error: %v", err)
	}
	defer store.Close()
	service, err = application.NewMemberService(dbsqlite.NewMemberRepository(store), gymID, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewMemberService after reopen returned error: %v", err)
	}
	members, err := NewApp(updater.NewClient("test"), "test", service).ListMembers()
	if err != nil {
		t.Fatalf("ListMembers after reopen returned error: %v", err)
	}
	if len(members) != 1 || members[0].ID != created.ID || members[0].FirstName != "Ada" {
		t.Fatalf("ListMembers after reopen = %#v, want Ada", members)
	}
}
