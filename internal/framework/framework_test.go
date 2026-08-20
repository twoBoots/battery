package framework

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListTemplates(t *testing.T) {
	templates := ListTemplates()
	if len(templates) == 0 {
		t.Fatalf("expected embedded templates to be listed, got empty slice")
	}

	foundRFC := false
	foundReview := false
	foundBatteryDoc := false

	for _, tmpl := range templates {
		if tmpl.Name == "skills/cooper-rfc" {
			foundRFC = true
		}
		if tmpl.Name == "skills/cooper-review" {
			foundReview = true
		}
		if tmpl.Name == "docs/BATTERY.md" {
			foundBatteryDoc = true
		}
	}

	if !foundRFC {
		t.Errorf("expected template 'skills/cooper-rfc' in ListTemplates")
	}
	if !foundReview {
		t.Errorf("expected template 'skills/cooper-review' in ListTemplates")
	}
	if !foundBatteryDoc {
		t.Errorf("expected template 'docs/BATTERY.md' in ListTemplates")
	}
}

func TestGetTemplate(t *testing.T) {
	t.Run("valid template", func(t *testing.T) {
		content, err := GetTemplate("skills/cooper-rfc")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(content, "cooper-rfc") {
			t.Errorf("expected content to contain 'cooper-rfc', got %q", content)
		}
	})

	t.Run("invalid template", func(t *testing.T) {
		_, err := GetTemplate("non-existent-template")
		if err == nil {
			t.Fatalf("expected error for non-existent template, got nil")
		}
	})
}

func TestInspectFrameworkStatus(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Initially empty directory -> all templates should be missing
	rep, err := InspectFrameworkStatus(tmpDir, "", "1.3.0")
	if err != nil {
		t.Fatalf("InspectFrameworkStatus failed: %v", err)
	}

	if rep.CLIVersion != "1.3.0" {
		t.Errorf("expected CLIVersion '1.3.0', got %q", rep.CLIVersion)
	}
	if rep.UpToDate {
		t.Errorf("expected UpToDate to be false for empty directory")
	}
	if !rep.UpdateAvailable {
		t.Errorf("expected UpdateAvailable to be true when templates are missing")
	}

	hasMissing := false
	for _, f := range rep.Files {
		if f.Status == StatusMissing {
			hasMissing = true
			break
		}
	}
	if !hasMissing {
		t.Errorf("expected some files to be missing in empty dir")
	}

	// 2. Populate some files with exact upstream content
	rfcContent, _ := GetTemplate("skills/cooper-rfc")
	rfcPath := filepath.Join(tmpDir, ".agents", "skills", "cooper-rfc", "SKILL.md")
	_ = os.MkdirAll(filepath.Dir(rfcPath), 0755)
	_ = os.WriteFile(rfcPath, []byte(rfcContent), 0644)

	// Populate a customized file
	reviewPath := filepath.Join(tmpDir, ".agents", "skills", "cooper-review", "SKILL.md")
	_ = os.MkdirAll(filepath.Dir(reviewPath), 0755)
	_ = os.WriteFile(reviewPath, []byte("# Custom team review skill with extra steps\n"), 0644)

	rep2, err := InspectFrameworkStatus(tmpDir, "", "1.3.0")
	if err != nil {
		t.Fatalf("InspectFrameworkStatus failed: %v", err)
	}

	var rfcStatus, reviewStatus *FileStatus
	for i := range rep2.Files {
		f := &rep2.Files[i]
		if f.TemplateName == "skills/cooper-rfc" {
			rfcStatus = f
		}
		if f.TemplateName == "skills/cooper-review" {
			reviewStatus = f
		}
	}

	if rfcStatus == nil || rfcStatus.Status != StatusUpToDate {
		t.Errorf("expected cooper-rfc to be UpToDate, got %+v", rfcStatus)
	}

	if reviewStatus == nil || reviewStatus.Status != StatusCustomizedLocally {
		t.Errorf("expected cooper-review to be CustomizedLocally, got %+v", reviewStatus)
	}
	if reviewStatus != nil && !reviewStatus.HasLocalModifications {
		t.Errorf("expected cooper-review HasLocalModifications to be true")
	}

	// 3. Test relative barrel path
	barrelDir := filepath.Join(tmpDir, "barrels", "auth")
	_ = os.MkdirAll(barrelDir, 0755)
	barrelRep, err := InspectFrameworkStatus(tmpDir, filepath.Join("barrels", "auth"), "1.3.0")
	if err != nil {
		t.Fatalf("InspectFrameworkStatus on barrel failed: %v", err)
	}
	if barrelRep.Target != filepath.Join("barrels", "auth") {
		t.Errorf("expected target %q, got %q", filepath.Join("barrels", "auth"), barrelRep.Target)
	}

	// 4. Test fully up to date workspace
	allUpToDateDir := t.TempDir()
	for _, tmpl := range ListTemplates() {
		content, _ := GetTemplate(tmpl.Name)
		targetFile := filepath.Join(allUpToDateDir, tmpl.TargetPath)
		_ = os.MkdirAll(filepath.Dir(targetFile), 0755)
		_ = os.WriteFile(targetFile, []byte(content), 0644)
	}

	allRep, err := InspectFrameworkStatus(allUpToDateDir, "", "1.3.0")
	if err != nil {
		t.Fatalf("InspectFrameworkStatus failed: %v", err)
	}
	if !allRep.UpToDate {
		t.Errorf("expected UpToDate to be true for fully populated dir")
	}
	if allRep.UpdateAvailable {
		t.Errorf("expected UpdateAvailable to be false for fully up-to-date dir")
	}
}
