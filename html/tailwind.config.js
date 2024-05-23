module.exports = {
  content: ["./templates/**/*.{html,js}"],
  safelist: [
    {pattern: /badge-+/},
  ],
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    darkTheme: "forest",
    themes: ["corporate", "forest"],
  },
};
