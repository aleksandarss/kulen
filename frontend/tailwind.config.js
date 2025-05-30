/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        primary: '#6B4F4F',   // muted brown
        secondary: '#DAD2BC', // parchment beige
        accent: '#A26769',    // warm clay
        background: '#F6F4F1' // paper white
      }
    }
  },
  plugins: [],
}
