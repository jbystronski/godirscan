package config

func (c *Config) ChangeTheme(t *Theme) {
	num := c.CurrentSchema
	if num < uint(len(c.ColorSchemas)-1) {
		num++
	} else {
		num = 0
	}

	c.CurrentSchema = num

	SetTheme(c.ColorSchemas[num], t)

	c.encode(c)
}

func SetTheme(schema Schema, theme *Theme) {
	/*

		setting foreground valules

	*/

	theme.Main = foreground[schema.Main]
	theme.Accent = foreground[schema.Accent]
	theme.Highlight = foreground[schema.Highlight]
	theme.Select = foreground[schema.Select]
	theme.Prompt = foreground[schema.Prompt]
	theme.Header = foreground[schema.Header]

	/*

		setting background values

	*/

	theme.BgHighlight = background[schema.BgHighlight]
	theme.BgHeader = background[schema.BgHeader]
	theme.BgSelect = background[schema.BgSelect]
	theme.BgPrompt = background[schema.BgPrompt]
}
