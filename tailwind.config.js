/** @type {import('tailwindcss').Config} */
//const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
	content: ["pub/**/*.{templ,go}"],
	theme: {
		fontFamily: {
			garamond: ['"adobe-garamond-pro"', 'serif'],
			mono: ['monospace'],
		},
		extend: {
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
		require("daisyui"),
	],
	daisyui: {
		themes: [
			"lofi"
		]
	},
	darkMode: "false",
}
