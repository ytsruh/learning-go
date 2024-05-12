import presetQuick from "franken-ui/shadcn-ui/preset-quick";

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.{templ,tmpl}"],
  theme: {
    extend: {},
  },
  plugins: [],
  presets: [presetQuick()],
};
