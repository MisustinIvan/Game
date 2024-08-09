package main

import (
	"os"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TextureManager struct {
	resources map[string]*ebiten.Image
	unknown   *ebiten.Image
}

func NewTextureManager(unknown_path string) (*TextureManager, error) {
	uknown, _, err := ebitenutil.NewImageFromFile(unknown_path)
	return &TextureManager{
		resources: make(map[string]*ebiten.Image),
		unknown:   uknown,
	}, err
}

func (r *TextureManager) LoadTexture(key string, path string) error {
	tex, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return err
	} else {
		r.resources[key] = tex
		return nil
	}
}

func (r *TextureManager) LoadTextures(dir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	dir = path.Join(cwd, dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(name, ".png") {
				key := strings.Split(name, ".")[0]
				err := r.LoadTexture(key, path.Join(dir, name))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *TextureManager) AddTexture(key string, tex *ebiten.Image) {
	r.resources[key] = tex
}

func (r TextureManager) GetTexture(key string) *ebiten.Image {
	tex, ok := r.resources[key]
	if !ok {
		return r.unknown
	} else {
		return tex
	}
}
