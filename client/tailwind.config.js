/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        'brand': {
          'primary': '#007bff',
          'secondary': '#6c757d',
          'success': '#28a745',
          'danger': '#dc3545',
          'warning': '#ffc107'
        },
        'background': {
          'light': '#f8f9fa',
          'dark': '#343a40'
        }
      },
      fontFamily: {
        'sans': ['Inter', 'system-ui', 'sans-serif'],
        'mono': ['Fira Code', 'monospace']
      },
    }
  },
  plugins: []
}