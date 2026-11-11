package input_test

import (
	"testing"
	"webtyp.com/input"
	"webtyp.com/model"
)

var (
	pngBytes = []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A, 0, 0}
	pdf = []byte("%PDF-1.7\x00\x00")
	mp4 = []byte{0, 0, 0, 0x18, 'f', 't', 'y', 'p', 'i', 's', 'o', 'm', 0, 0}
	wav = []byte("RIFF\x24\x00\x00\x00WAVEfmt \x00\x00")
)

func TestImageFile_AcceptsImageRejectsPDF(t *testing.T) {
	k := input.ImageFile(1024)
	if _, err := k.ValidateFile(pngBytes); err != nil {
		t.Fatalf("expected no error for png, got %v", err)
	}
	if _, err := k.ValidateFile(pdf); err == nil {
		t.Errorf("expected error for pdf, got nil")
	}
}

func TestDocumentFile_AcceptsPDFRejectsImage(t *testing.T) {
	k := input.DocumentFile(1024)
	if _, err := k.ValidateFile(pdf); err != nil {
		t.Fatalf("expected no error for pdf, got %v", err)
	}
	if _, err := k.ValidateFile(pngBytes); err == nil {
		t.Errorf("expected error for png, got nil")
	}
}

func TestScanFile_AcceptsImageAndPDF(t *testing.T) {
	k := input.ScanFile(1024)
	if _, err := k.ValidateFile(pngBytes); err != nil {
		t.Fatalf("expected no error for png, got %v", err)
	}
	if _, err := k.ValidateFile(pdf); err != nil {
		t.Fatalf("expected no error for pdf, got %v", err)
	}
}

func TestVideoFile_AcceptsMP4RejectsAudio(t *testing.T) {
	k := input.VideoFile(1024)
	if _, err := k.ValidateFile(mp4); err != nil {
		t.Fatalf("expected no error for mp4, got %v", err)
	}
	if _, err := k.ValidateFile(wav); err == nil {
		t.Errorf("expected error for wav, got nil")
	}
}

func TestAudioFile_AcceptsWAVRejectsVideo(t *testing.T) {
	k := input.AudioFile(1024)
	if _, err := k.ValidateFile(wav); err != nil {
		t.Fatalf("expected no error for wav, got %v", err)
	}
	if _, err := k.ValidateFile(mp4); err == nil {
		t.Errorf("expected error for mp4, got nil")
	}
}

func TestFileKinds_StoreText(t *testing.T) {
	kinds := []input.FileKind{
		input.ImageFile(10),
		input.DocumentFile(20),
		input.ScanFile(30),
		input.VideoFile(40),
		input.AudioFile(50),
	}
	for i, k := range kinds {
		if k.Storage() != model.FieldText {
			t.Errorf("kind %d expected storage FieldText, got %v", i, k.Storage())
		}
		if k.Name() != "file" {
			t.Errorf("kind %d expected Name file, got %v", i, k.Name())
		}
		if k.MaxBytes() != (i+1)*10 {
			t.Errorf("kind %d expected MaxBytes %d, got %v", i, (i+1)*10, k.MaxBytes())
		}
	}
}
