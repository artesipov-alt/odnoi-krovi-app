console.time("Build time");
console.log("🔨 Building project...");

await Bun.build({
  entrypoints: ["./bot.ts"],
  outdir: "./dist",
  target: "bun",
  minify: true,
  sourcemap: true,
  packages: "external",
  format: "esm",
});

console.log("✅ Build completed successfully!");
console.timeEnd("Build time");
