/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors:{
        primary: "var(--clr-primary)",
        secondary: "var(--clr-secondary)",
      }
    },
  },
  plugins: [],
}

