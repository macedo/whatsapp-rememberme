/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/views/layouts/*.html",
    "./web/views/partials/*.html",
    "./web/views/*.html"
  ],
  theme: {
    extend: {},
  },
  plugins: [
      require('@tailwindcss/forms'),
  ],
}
