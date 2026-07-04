/**
 * Сборка max-bot для production.
 *
 * Использует esbuild (вместо `tsc`) по двум причинам:
 *   1. Поддерживает `verbatimModuleSyntax: true` в tsconfig — не добавляет
 *      "use strict", не переписывает ESM-импорты.
 *   2. Не требует явных `.js` расширений в импортах (NodeNext этого
 *      требует, esbuild — нет).
 *
 * Результат: один ESM-бандл в `dist/bot.js` со всеми зависимостями
 * (как в Bun.build, но под Node.js).
 */
import { build } from "esbuild";
import { rmSync } from "node:fs";

console.log(" Building max-bot...");

const start = Date.now();

// Чистим старый dist
rmSync("./dist", { recursive: true, force: true });

await build({
  entryPoints: ["./bot.ts"],
  outfile: "./dist/bot.js",
  bundle: true,
  platform: "node",
  target: "node24",
  format: "esm",
  packages: "external", // node_modules в bundle не тащим
  minify: true,
  sourcemap: true,
  banner: {
    // Подсказка для Node, что это ESM
    js: "import { createRequire as __cR } from 'node:module'; const require = __cR(import.meta.url);",
  },
  // shared/ts лежит уровнем выше, esbuild сам разберётся
  absWorkingDir: process.cwd(),
  logLevel: "info",
});

console.log(`Build completed in ${Date.now() - start}ms`);
