package tool

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"tool/resource/fonts"
)

type MyTheme struct {
	regular, bold, italic, boldItalic, monospace fyne.Resource
}

func (t *MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t *MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return m.monospace
	}
	if style.Bold {
		if style.Italic {
			return m.boldItalic
		}
		return m.bold
	}
	if style.Italic {
		return m.italic
	}
	return m.regular
}

func (m *MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (t *MyTheme) SetFonts(regularFontPath string, monoFontPath string) {
	t.regular = theme.TextFont()
	t.bold = theme.TextBoldFont()
	t.italic = theme.TextItalicFont()
	t.boldItalic = theme.TextBoldItalicFont()
	t.monospace = theme.TextMonospaceFont()

	if regularFontPath == "bindata" {
		t.regular = loadCustomFont(regularFontPath, "Regular", t.regular)
		t.bold = loadCustomFont(regularFontPath, "Bold", t.bold)
		t.italic = loadCustomFont(regularFontPath, "Italic", t.italic)
		t.boldItalic = loadCustomFont(regularFontPath, "BoldItalic", t.boldItalic)
	}
	if monoFontPath != "" {
		t.monospace = loadCustomFont(monoFontPath, "Regular", t.monospace)
	} else {
		t.monospace = t.regular
	}
}

func loadCustomFont(env, variant string, fallback fyne.Resource) fyne.Resource {
	b, _ := fonts.Asset(fonts.AssetNames()[0])

	return fyne.NewStaticResource("font", b)
}
