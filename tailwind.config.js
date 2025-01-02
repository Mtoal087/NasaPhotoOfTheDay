/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./templates/*.html'],
  theme: {
    extend: {
      transitionProperty: {
        'size': 'width, height'
      }
    },
  },
  plugins: [],
}

