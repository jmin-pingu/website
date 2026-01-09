/** @type {import('tailwindcss').Config} */
//const defaultTheme = require('tailwindcss/defaultTheme')
//

module.exports = {
	content: ["./internal/pub/**/*.{templ,go}", "./**/*.{templ,go}"],
	important: true,
	theme: {
		extend: {
			cursor: {
				'creative-default': 'url(/assets/cursors/normal_select_1.cur), default',
				'creative-text': 'url(/assets/cursors/text_select_1.cur), default',
				'creative-pointer': 'url(/assets/cursors/link_select.cur), default',
        		'custom': 'url(https://play.tailwindcss.com/favicons/favicon-16x16.png?v=3), default',
			},
			colors: {
				gray: {
					350: "#b8bdc6",
				}
			},
			fontFamily: {
				tnr: ['"Times New Roman"', 'serif'],
				garamond: ['"adobe-garamond-pro"', 'serif'],
				mono: ['monospace'],
			},
			spacing: {
				'128': '32rem',
				'144': '36rem',
			},
			borderRadius: {
				'4xl': '2rem',
			}
		}
	},
	plugins: [
		require("@tailwindcss/typography"),
	],
	daisyui: {
		themes: [
			"lofi"
		]
	},
	darkMode: "false",
}
