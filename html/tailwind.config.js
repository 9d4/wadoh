module.exports = {
  content: ["./templates/**/*.{html,js}"],
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    darkTheme: "forest",
    themes: ["corporate", "forest"],
  },
};
