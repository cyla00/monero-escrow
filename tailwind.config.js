module.exports = {
  content: ["./views/**/*.{html,js,go,templ}", "./components/**/*.{html,js,go,templ}"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]
}