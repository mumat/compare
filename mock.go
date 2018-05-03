package main

type mockAssets string

func (assets mockAssets) String(name string) string {
	return string(assets)
}
