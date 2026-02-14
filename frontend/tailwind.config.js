/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        background: "#0F1115",
        surface: "#161A22",
        border: "#262B36",
        primary: "#E6EAF2",
        muted: "#8B93A7",
        accent: {
          DEFAULT: "#4F7CFF",
          hover: "#3D63D9",
        },
      },
      borderRadius: {
        sm: "2px",
        DEFAULT: "4px",
      },
      fontFamily: {
        mono: ["'JetBrains Mono'", "monospace"],
        sans: ["'Inter'", "system-ui", "sans-serif"],
      },
    },
  },
  plugins: [],
};
