package handler

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestHandleConvertEncoding_UTF8ToCP1251(t *testing.T) {
	tempDir := t.TempDir()
	h := NewHandler([]string{tempDir})

	testFile := filepath.Join(tempDir, "test.txt")
	// UTF-8 content with Cyrillic
	utf8Content := "Привет мир" // "Hello world" in Russian
	os.WriteFile(testFile, []byte(utf8Content), 0644)

	result, output, err := h.HandleConvertEncoding(context.Background(), nil, ConvertEncodingInput{
		Path: testFile,
		From: "utf-8",
		To:   "cp1251",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Errorf("expected success, got error")
	}
	if output.SourceEncoding != "utf-8" {
		t.Errorf("expected source encoding utf-8, got %s", output.SourceEncoding)
	}
	if output.TargetEncoding != "cp1251" {
		t.Errorf("expected target encoding cp1251, got %s", output.TargetEncoding)
	}

	// Verify file was converted (CP1251 bytes are different from UTF-8)
	converted, _ := os.ReadFile(testFile)
	if string(converted) == utf8Content {
		t.Error("file content should be different after conversion")
	}
}

func TestHandleConvertEncoding_CP1251ToUTF8(t *testing.T) {
	tempDir := t.TempDir()
	h := NewHandler([]string{tempDir})

	testFile := filepath.Join(tempDir, "test.txt")
	// CP1251 bytes for "Привет" (Russian "Hello")
	cp1251Bytes := []byte{0xCF, 0xF0, 0xE8, 0xE2, 0xE5, 0xF2}
	os.WriteFile(testFile, cp1251Bytes, 0644)

	result, output, err := h.HandleConvertEncoding(context.Background(), nil, ConvertEncodingInput{
		Path: testFile,
		From: "cp1251",
		To:   "utf-8",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Errorf("expected success, got error")
	}
	if output.TargetEncoding != "utf-8" {
		t.Errorf("expected target encoding utf-8, got %s", output.TargetEncoding)
	}

	// Verify file is now valid UTF-8
	converted, _ := os.ReadFile(testFile)
	expected := "Привет"
	if string(converted) != expected {
		t.Errorf("expected %q, got %q", expected, string(converted))
	}
}

func TestHandleConvertEncoding_WithBackup(t *testing.T) {
	tempDir := t.TempDir()
	h := NewHandler([]string{tempDir})

	testFile := filepath.Join(tempDir, "test.txt")
	originalContent := []byte("original content")
	os.WriteFile(testFile, originalContent, 0644)

	result, output, err := h.HandleConvertEncoding(context.Background(), nil, ConvertEncodingInput{
		Path:   testFile,
		From:   "utf-8",
		To:     "cp1251",
		Backup: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Errorf("expected success, got error")
	}
	if output.BackupPath == "" {
		t.Error("expected backup path to be set")
	}

	// Verify backup file exists with original content
	backupContent, err := os.ReadFile(output.BackupPath)
	if err != nil {
		t.Errorf("backup file should exist: %v", err)
	}
	if string(backupContent) != string(originalContent) {
		t.Error("backup should contain original content")
	}
}

func TestHandleConvertEncoding_MissingTo(t *testing.T) {
	tempDir := t.TempDir()
	h := NewHandler([]string{tempDir})

	testFile := filepath.Join(tempDir, "test.txt")
	os.WriteFile(testFile, []byte("content"), 0644)

	result, _, err := h.HandleConvertEncoding(context.Background(), nil, ConvertEncodingInput{
		Path: testFile,
		From: "utf-8",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Error("expected error for missing 'to' parameter")
	}
}

func TestHandleConvertEncoding_OutsideAllowed(t *testing.T) {
	tempDir := t.TempDir()
	h := NewHandler([]string{tempDir})

	result, _, err := h.HandleConvertEncoding(context.Background(), nil, ConvertEncodingInput{
		Path: "/some/random/file.txt",
		To:   "utf-8",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Error("expected error for path outside allowed directories")
	}
}

func TestHandleConvertEncoding_UTF8ToGBK(t *testing.T) {
	tempDir := t.TempDir()
	h := NewHandler([]string{tempDir})

	testFile := filepath.Join(tempDir, "test.txt")
	// UTF-8 content with Chinese characters
	utf8Content := "你好世界" // "Hello world" in Chinese
	os.WriteFile(testFile, []byte(utf8Content), 0644)

	result, output, err := h.HandleConvertEncoding(context.Background(), nil, ConvertEncodingInput{
		Path: testFile,
		From: "utf-8",
		To:   "gbk",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Errorf("expected success, got error")
	}
	if output.SourceEncoding != "utf-8" {
		t.Errorf("expected source encoding utf-8, got %s", output.SourceEncoding)
	}
	if output.TargetEncoding != "gbk" {
		t.Errorf("expected target encoding gbk, got %s", output.TargetEncoding)
	}

	// Verify file was converted (GBK bytes are different from UTF-8)
	converted, _ := os.ReadFile(testFile)
	if string(converted) == utf8Content {
		t.Error("file content should be different after conversion")
	}
}

func TestHandleConvertEncoding_GBKToUTF8(t *testing.T) {
	tempDir := t.TempDir()
	h := NewHandler([]string{tempDir})

	testFile := filepath.Join(tempDir, "test.txt")
	// GBK bytes for "你好" (Chinese "Hello")
	gbkBytes := []byte{0xC4, 0xE3, 0xBA, 0xC3}
	os.WriteFile(testFile, gbkBytes, 0644)

	result, output, err := h.HandleConvertEncoding(context.Background(), nil, ConvertEncodingInput{
		Path: testFile,
		From: "gbk",
		To:   "utf-8",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Errorf("expected success, got error")
	}
	if output.TargetEncoding != "utf-8" {
		t.Errorf("expected target encoding utf-8, got %s", output.TargetEncoding)
	}

	// Verify file is now valid UTF-8
	converted, _ := os.ReadFile(testFile)
	expected := "你好"
	if string(converted) != expected {
		t.Errorf("expected %q, got %q", expected, string(converted))
	}
}

func TestHandleConvertEncoding_GB18030ToUTF8(t *testing.T) {
	tempDir := t.TempDir()
	h := NewHandler([]string{tempDir})

	testFile := filepath.Join(tempDir, "test.txt")
	// GB18030 bytes for "你好世界" (Chinese "Hello world")
	// Same as GBK for common characters
	gb18030Bytes := []byte{0xC4, 0xE3, 0xBA, 0xC3, 0xCA, 0xC0, 0xBD, 0xE7}
	os.WriteFile(testFile, gb18030Bytes, 0644)

	result, output, err := h.HandleConvertEncoding(context.Background(), nil, ConvertEncodingInput{
		Path: testFile,
		From: "gb18030",
		To:   "utf-8",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Errorf("expected success, got error")
	}
	if output.TargetEncoding != "utf-8" {
		t.Errorf("expected target encoding utf-8, got %s", output.TargetEncoding)
	}

	// Verify file is now valid UTF-8
	converted, _ := os.ReadFile(testFile)
	expected := "你好世界"
	if string(converted) != expected {
		t.Errorf("expected %q, got %q", expected, string(converted))
	}
}

func TestHandleConvertEncoding_GB2312Alias(t *testing.T) {
	tempDir := t.TempDir()
	h := NewHandler([]string{tempDir})

	testFile := filepath.Join(tempDir, "test.txt")
	// GBK/GB2312 bytes for "测试" (Chinese "test")
	gb2312Bytes := []byte{0xB2, 0xE2, 0xCA, 0xD4}
	os.WriteFile(testFile, gb2312Bytes, 0644)

	result, output, err := h.HandleConvertEncoding(context.Background(), nil, ConvertEncodingInput{
		Path: testFile,
		From: "gb2312",
		To:   "utf-8",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Errorf("expected success, got error (gb2312 should work as alias for gbk)")
	}
	if output.SourceEncoding != "gb2312" {
		t.Errorf("expected source encoding gb2312, got %s", output.SourceEncoding)
	}

	// Verify file is now valid UTF-8
	converted, _ := os.ReadFile(testFile)
	expected := "测试"
	if string(converted) != expected {
		t.Errorf("expected %q, got %q", expected, string(converted))
	}
}
