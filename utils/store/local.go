/*
 * Copyright (C) 2020-2023 Arm Limited or its affiliates and Contributors. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package store

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/ARM-software/golang-utils/utils/commonerrors"
	"github.com/ARM-software/golang-utils/utils/filesystem"
	"github.com/ARM-software/golang-utils/utils/reflection"
)

type LocalStore struct {
	mu       sync.RWMutex
	rootPath string
	fs       filesystem.FS
}

func (s *LocalStore) GetFilesystem() filesystem.FS {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.fs
}

func (s *LocalStore) Exists() bool {
	if reflection.IsEmpty(s.GetPath()) {
		return false
	}
	return s.fs.Exists(s.GetPath())
}

func (s *LocalStore) Close() error {
	return nil
}

func (s *LocalStore) GetPath() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rootPath
}

func (s *LocalStore) SetPath(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rootPath = path
	return nil
}

func (s *LocalStore) Create(context.Context) error {
	return s.fs.MkDir(s.GetPath())
}

func (s *LocalStore) Clear(ctx context.Context) error {
	if !s.Exists() {
		return nil
	}
	return s.fs.CleanDirWithContext(ctx, s.GetPath())
}

func (s *LocalStore) GetElementPath(elementName string) (path string) {
	if !s.Exists() {
		return
	}
	if !filepath.IsAbs(elementName) {
		path = filepath.Clean(elementName)
		return
	}
	paths, err := s.fs.ConvertToRelativePath(s.GetPath(), elementName)
	if err != nil || len(paths) < 1 {
		fmt.Println(paths, err)
		return
	}
	path = paths[0]
	return
}

func (s *LocalStore) HasElement(elementName string) bool {
	if !s.Exists() {
		return false
	}
	return s.fs.Exists(filepath.Join(s.GetPath(), s.GetElementPath(elementName)))
}

// NewLocalStore creates a persistent local store in the filesystem.
func NewLocalStore(rootPath string) IStore {
	return &LocalStore{
		mu:       sync.RWMutex{},
		rootPath: rootPath,
		fs:       filesystem.GetGlobalFileSystem(),
	}
}

type LocalTemporaryStore struct {
	LocalStore
}

func (s *LocalTemporaryStore) Close() error {
	return s.fs.Rm(s.GetPath())
}

func (s *LocalTemporaryStore) SetPath(string) error {
	return commonerrors.ErrUnsupported
}
func (s *LocalTemporaryStore) Create(context.Context) error {
	name, err := s.fs.TempDirInTempDir("local-store-")
	if err != nil {
		return err
	}
	return s.LocalStore.SetPath(name)
}

// NewLocalTemporaryStore creates a local store in the temp directory
func NewLocalTemporaryStore() IStore {
	return &LocalTemporaryStore{
		LocalStore: LocalStore{
			mu:       sync.RWMutex{},
			rootPath: "",
			fs:       filesystem.GetGlobalFileSystem(),
		},
	}
}
