module.exports = {
    content: ["./views/**/*.{templ,go,html}", "./components/**/*.{templ,go,html}"],
    theme: {
        extend: {},
    },
    plugins: [
        require('@tailwindcss/forms'),
        require('@tailwindcss/typography'),
    ],
    corePlugins: {
        preflight: true,
    }
}