'use strict';

Object.defineProperty(exports, '__esModule', { value: true });

var transformCss = require('@vanilla-extract/css/transformCss');
var evalCode = require('eval');
var javascriptStringify = require('javascript-stringify');
var dedent = require('dedent');
var crypto = require('crypto');
var zlib = require('zlib');
var util = require('util');
var path = require('path');
var findUp = require('find-up');
var fs = require('fs');
var esbuild = require('esbuild');
var babel = require('@babel/core');
var vanillaBabelPlugin = require('@vanilla-extract/babel-plugin-debug-ids');
var typescriptSyntax = require('@babel/plugin-syntax-typescript');
var mlly = require('mlly');

function _interopDefault (e) { return e && e.__esModule ? e : { 'default': e }; }

function _interopNamespace(e) {
  if (e && e.__esModule) return e;
  var n = Object.create(null);
  if (e) {
    Object.keys(e).forEach(function (k) {
      if (k !== 'default') {
        var d = Object.getOwnPropertyDescriptor(e, k);
        Object.defineProperty(n, k, d.get ? d : {
          enumerable: true,
          get: function () { return e[k]; }
        });
      }
    });
  }
  n["default"] = e;
  return Object.freeze(n);
}

var evalCode__default = /*#__PURE__*/_interopDefault(evalCode);
var dedent__default = /*#__PURE__*/_interopDefault(dedent);
var crypto__default = /*#__PURE__*/_interopDefault(crypto);
var path__default = /*#__PURE__*/_interopDefault(path);
var findUp__default = /*#__PURE__*/_interopDefault(findUp);
var babel__namespace = /*#__PURE__*/_interopNamespace(babel);
var vanillaBabelPlugin__default = /*#__PURE__*/_interopDefault(vanillaBabelPlugin);
var typescriptSyntax__default = /*#__PURE__*/_interopDefault(typescriptSyntax);

const hash = value => crypto__default["default"].createHash('md5').update(value).digest('hex');

const zip = util.promisify(zlib.gzip);
const unzip = util.promisify(zlib.gunzip);

// The byte threshold for applying compression, below which compressing would out-weigh its value.
const compressionThreshold = 1000;
const compressionFlag = '#';
async function serializeCss(source) {
  if (source.length > compressionThreshold) {
    const compressedSource = await zip(source);
    return compressionFlag + compressedSource.toString('base64');
  }
  return Buffer.from(source, 'utf-8').toString('base64');
}
async function deserializeCss(source) {
  if (source.indexOf(compressionFlag) > -1) {
    const decompressedSource = await unzip(Buffer.from(source.replace(compressionFlag, ''), 'base64'));
    return decompressedSource.toString('utf-8');
  }
  return Buffer.from(source, 'base64').toString('utf-8');
}

const originalNodeEnv =         "production";

// Copied from https://github.com/you-dont-need/You-Dont-Need-Lodash-Underscore/blob/51f83bd3db728fd7ee177de1ffc253fdb99c537f/README.md#_isplainobject
function isPlainObject(value) {
  if (typeof value !== 'object' || value === null) {
    return false;
  }
  if (Object.prototype.toString.call(value) !== '[object Object]') {
    return false;
  }
  const prototype = Object.getPrototypeOf(value);
  if (prototype === null) {
    return true;
  }
  const constructor = Object.prototype.hasOwnProperty.call(prototype, 'constructor') && prototype.constructor;
  return typeof constructor === 'function' && constructor instanceof constructor && Function.prototype.call(constructor) === Function.prototype.call(value);
}
function stringifyFileScope({
  packageName,
  filePath
}) {
  return packageName ? `${filePath}$$$${packageName}` : filePath;
}
function parseFileScope(serialisedFileScope) {
  const [filePath, packageName] = serialisedFileScope.split('$$$');
  return {
    filePath,
    packageName
  };
}
async function processVanillaFile({
  source,
  filePath,
  outputCss = true,
  identOption = 'short' ,
  serializeVirtualCssPath
}) {
  const cssByFileScope = new Map();
  const localClassNames = new Set();
  const composedClassLists = [];
  const usedCompositions = new Set();
  const cssAdapter = {
    appendCss: (css, fileScope) => {
      if (outputCss) {
        const serialisedFileScope = stringifyFileScope(fileScope);
        const fileScopeCss = cssByFileScope.get(serialisedFileScope) ?? [];
        fileScopeCss.push(css);
        cssByFileScope.set(serialisedFileScope, fileScopeCss);
      }
    },
    registerClassName: className => {
      localClassNames.add(className);
    },
    registerComposition: composedClassList => {
      composedClassLists.push(composedClassList);
    },
    markCompositionUsed: identifier => {
      usedCompositions.add(identifier);
    },
    onEndFileScope: () => {},
    getIdentOption: () => identOption
  };
  const currentNodeEnv =         "production";

  // Vite sometimes modifies NODE_ENV which causes different versions (e.g. dev/prod) of vanilla packages to be loaded
  // This can cause CSS to be bound to the wrong instance, resulting in no CSS output
  // To get around this we set the NODE_ENV back to the original value ONLY during eval
  process.env.NODE_ENV = originalNodeEnv;
  const adapterBoundSource = `
    require('@vanilla-extract/css/adapter').setAdapter(__adapter__);
    ${source}
  `;
  const evalResult = evalCode__default["default"](adapterBoundSource, filePath, {
    console,
    process,
    __adapter__: cssAdapter
  }, true);
  process.env.NODE_ENV = currentNodeEnv;
  const cssImports = [];
  for (const [serialisedFileScope, fileScopeCss] of cssByFileScope) {
    const fileScope = parseFileScope(serialisedFileScope);
    const css = transformCss.transformCss({
      localClassNames: Array.from(localClassNames),
      composedClassLists,
      cssObjs: fileScopeCss
    }).join('\n');
    const fileName = `${fileScope.filePath}.vanilla.css`;
    let virtualCssFilePath;
    if (serializeVirtualCssPath) {
      const serializedResult = serializeVirtualCssPath({
        fileName,
        fileScope,
        source: css
      });
      if (typeof serializedResult === 'string') {
        virtualCssFilePath = serializedResult;
      } else {
        virtualCssFilePath = await serializedResult;
      }
    } else {
      const serializedCss = await serializeCss(css);
      virtualCssFilePath = `import '${fileName}?source=${serializedCss}';`;
    }
    cssImports.push(virtualCssFilePath);
  }

  // We run this code inside eval as jest seems to create a difrerent instance of the adapter file
  // for requires executed within the eval and all CSS can be lost.
  evalCode__default["default"](`const { removeAdapter } = require('@vanilla-extract/css/adapter');
    // Backwards compat with older versions of @vanilla-extract/css
    if (removeAdapter) {
      removeAdapter();
    }
  `, filePath, {
    console,
    process
  }, true);
  const unusedCompositions = composedClassLists.filter(({
    identifier
  }) => !usedCompositions.has(identifier)).map(({
    identifier
  }) => identifier);
  const unusedCompositionRegex = unusedCompositions.length > 0 ? RegExp(`(${unusedCompositions.join('|')})\\s`, 'g') : null;
  return serializeVanillaModule(cssImports, evalResult, unusedCompositionRegex);
}
function stringifyExports(functionSerializationImports, value, unusedCompositionRegex, key, exportLookup, exportDependencyGraph) {
  return javascriptStringify.stringify(value, (value, _indent, next) => {
    const valueType = typeof value;
    if (valueType === 'boolean' || valueType === 'number' || valueType === 'undefined' || value === null) {
      return next(value);
    }
    if (Array.isArray(value) || isPlainObject(value)) {
      const reusedExport = exportLookup.get(value);
      if (reusedExport && reusedExport !== key) {
        exportDependencyGraph.addDependency(key, reusedExport);
        return reusedExport;
      }
      return next(value);
    }
    if (Symbol.toStringTag in Object(value)) {
      const {
        [Symbol.toStringTag]: _tag,
        ...valueWithoutTag
      } = value;
      return next(valueWithoutTag);
    }
    if (valueType === 'string') {
      return next(unusedCompositionRegex ? value.replace(unusedCompositionRegex, '') : value);
    }
    if (valueType === 'function' && (value.__function_serializer__ || value.__recipe__)) {
      const {
        importPath,
        importName,
        args
      } = value.__function_serializer__ || value.__recipe__;
      if (typeof importPath !== 'string' || typeof importName !== 'string' || !Array.isArray(args)) {
        throw new Error('Invalid function serialization params');
      }
      try {
        const hashedImportName = `_${hash(`${importName}${importPath}`).slice(0, 5)}`;
        functionSerializationImports.add(`import { ${importName} as ${hashedImportName} } from '${importPath}';`);
        return `${hashedImportName}(${args.map(arg => stringifyExports(functionSerializationImports, arg, unusedCompositionRegex, key, exportLookup, exportDependencyGraph)).join(',')})`;
      } catch (err) {
        console.error(err);
        throw new Error('Invalid function serialization params');
      }
    }
    throw new Error(dedent__default["default"]`
        Invalid exports.

        You can only export plain objects, arrays, strings, numbers and null/undefined.
      `);
  }, 0, {
    references: true,
    // Allow circular references
    maxDepth: Infinity,
    maxValues: Infinity
  });
}
const defaultExportName = '__default__';
class DependencyGraph {
  constructor() {
    this.graph = new Map();
  }

  /**
   * Creates a "depends on" relationship between `key` and `dependency`
   */
  addDependency(key, dependency) {
    const dependencies = this.graph.get(key);
    if (dependencies) {
      dependencies.add(dependency);
    } else {
      this.graph.set(key, new Set([dependency]));
    }
  }

  /**
   * Whether or not `key` depends on `dependency`
   */
  dependsOn(key, dependency) {
    const dependencies = this.graph.get(key);
    if (dependencies) {
      if (dependencies !== null && dependencies !== void 0 && dependencies.has(dependency)) {
        return true;
      }
      for (const [dep] of dependencies.entries()) {
        if (this.dependsOn(dep, dependency)) {
          return true;
        }
      }
    }
    return false;
  }
}
function serializeVanillaModule(cssImports, exports, unusedCompositionRegex) {
  const functionSerializationImports = new Set();
  const exportLookup = new Map(Object.entries(exports).map(([key, value]) => [value, key === 'default' ? defaultExportName : key]));
  const exportDependencyGraph = new DependencyGraph();
  const moduleExports = Object.entries(exports).map(([key, value]) => {
    const serializedExport = stringifyExports(functionSerializationImports, value, unusedCompositionRegex, key === 'default' ? defaultExportName : key, exportLookup, exportDependencyGraph);
    if (key === 'default') {
      return [defaultExportName, [`var ${defaultExportName} = ${serializedExport};`, `export default ${defaultExportName};`].join('\n')];
    }
    return [key, `export var ${key} = ${serializedExport};`];
  });
  const sortedModuleExports = moduleExports.sort(([key1], [key2]) => {
    if (exportDependencyGraph.dependsOn(key1, key2)) {
      return 1;
    }
    if (exportDependencyGraph.dependsOn(key2, key1)) {
      return -1;
    }
    return 0;
  }).map(([, s]) => s);
  const outputCode = [...cssImports, ...functionSerializationImports, ...sortedModuleExports];
  return outputCode.join('\n');
}

async function getSourceFromVirtualCssFile(id) {
  const match = id.match(/^(?<fileName>.*)\?source=(?<source>.*)$/);
  if (!match || !match.groups) {
    throw new Error('No source in vanilla CSS file');
  }
  const source = await deserializeCss(match.groups.source);
  return {
    fileName: match.groups.fileName,
    source
  };
}

function getClosestPackageInfo(directory) {
  const packageJsonPath = findUp__default["default"].sync('package.json', {
    cwd: directory
  });
  if (packageJsonPath) {
    const {
      name
    } = require(packageJsonPath);
    return {
      name,
      path: packageJsonPath,
      dirname: path__default["default"].dirname(packageJsonPath)
    };
  }
}
const packageInfoCache = new Map();
function getPackageInfo(cwd) {
  const resolvedCwd = cwd ?? process.cwd();
  const cachedValue = packageInfoCache.get(resolvedCwd);
  if (cachedValue) {
    return cachedValue;
  }
  let packageInfo = getClosestPackageInfo(resolvedCwd);
  while (packageInfo && !packageInfo.name) {
    packageInfo = getClosestPackageInfo(path__default["default"].resolve(packageInfo.dirname, '..'));
  }
  if (!packageInfo || !packageInfo.name) {
    throw new Error(`Couldn't find parent package.json with a name field from '${resolvedCwd}'`);
  }
  packageInfoCache.set(resolvedCwd, packageInfo);
  return packageInfo;
}

// Vite adds a "?used" to CSS files it detects, this isn't relevant for
// .css.ts files but it's added anyway so we need to allow for it in the file match
const cssFileFilter = /\.css\.(js|cjs|mjs|jsx|ts|tsx)(\?used)?$/;
const virtualCssFileFilter = /\.vanilla\.css\?source=.*$/;

// Inlined from @rollup/pluginutils
// https://github.com/rollup/plugins/blob/33174f956304ab4aad4bbaba656f627c31679dc5/packages/pluginutils/src/normalizePath.ts#L5-L7
const normalizePath = filename => filename.split(path.win32.sep).join(path.posix.sep);
function addFileScope({
  source,
  filePath,
  rootPath,
  packageName,
  globalAdapterIdentifier
}) {
  // Encode windows file paths as posix
  const normalizedPath = normalizePath(path.relative(rootPath, filePath));
  const {
    hasESM,
    isMixed
  } = mlly.detectSyntax(source);
  if (source.includes('@vanilla-extract/css/fileScope')) {
    source = source.replace(/setFileScope\(((\n|.)*?)\)/, `setFileScope("${normalizedPath}", "${packageName}")`);
  } else {
    if (hasESM && !isMixed) {
      source = dedent__default["default"](`
        import { setFileScope, endFileScope } from "@vanilla-extract/css/fileScope";
        setFileScope("${normalizedPath}", "${packageName}");
        ${source}
        endFileScope();
      `);
    } else {
      source = dedent__default["default"](`
        const __vanilla_filescope__ = require("@vanilla-extract/css/fileScope");
        __vanilla_filescope__.setFileScope("${normalizedPath}", "${packageName}");
        ${source}
        __vanilla_filescope__.endFileScope();
      `);
    }
  }
  if (globalAdapterIdentifier) {
    const adapterImport = hasESM && !isMixed ? 'import * as __vanilla_css_adapter__ from "@vanilla-extract/css/adapter";' : 'const __vanilla_css_adapter__ = require("@vanilla-extract/css/adapter");';
    source = dedent__default["default"](`
      ${adapterImport}
      __vanilla_css_adapter__.setAdapter(${globalAdapterIdentifier});
      ${source}
      __vanilla_css_adapter__.removeAdapter();
    `);
  }
  return source;
}

const transformSync = ({
  source,
  filePath,
  rootPath,
  packageName,
  identOption
}) => {
  let code = source;
  if (identOption === 'debug') {
    const result = babel__namespace.transformSync(source, {
      filename: filePath,
      cwd: rootPath,
      plugins: [vanillaBabelPlugin__default["default"], typescriptSyntax__default["default"]],
      configFile: false,
      babelrc: false
    });
    if (!result || result.code == null) {
      throw new Error('Error adding debug IDs');
    }
    code = result.code;
  }
  return addFileScope({
    source: code,
    filePath,
    rootPath,
    packageName
  });
};
const transform = async ({
  source,
  filePath,
  rootPath,
  packageName,
  identOption,
  globalAdapterIdentifier
}) => {
  let code = source;
  if (identOption === 'debug') {
    const result = await babel__namespace.transformAsync(source, {
      filename: filePath,
      cwd: rootPath,
      plugins: [vanillaBabelPlugin__default["default"], typescriptSyntax__default["default"]],
      configFile: false,
      babelrc: false
    });
    if (!result || result.code == null) {
      throw new Error('Error adding debug IDs');
    }
    code = result.code;
  }
  return addFileScope({
    source: code,
    filePath,
    rootPath,
    packageName,
    globalAdapterIdentifier
  });
};

const vanillaExtractTransformPlugin = ({
  identOption
}) => ({
  name: 'vanilla-extract-filescope',
  setup(build) {
    const packageInfo = getPackageInfo(build.initialOptions.absWorkingDir);
    build.onLoad({
      filter: cssFileFilter
    }, async ({
      path: path$1
    }) => {
      const originalSource = await fs.promises.readFile(path$1, 'utf-8');
      const source = await transform({
        source: originalSource,
        filePath: path$1,
        rootPath: build.initialOptions.absWorkingDir,
        packageName: packageInfo.name,
        identOption: identOption ?? (build.initialOptions.minify ? 'short' : 'debug')
      });
      return {
        contents: source,
        loader: path$1.match(/\.(ts|tsx)$/i) ? 'ts' : undefined,
        resolveDir: path.dirname(path$1)
      };
    });
  }
});
async function compile({
  filePath,
  identOption,
  cwd = process.cwd(),
  esbuildOptions
}) {
  const result = await esbuild.build({
    entryPoints: [filePath],
    metafile: true,
    bundle: true,
    external: ['@vanilla-extract', ...((esbuildOptions === null || esbuildOptions === void 0 ? void 0 : esbuildOptions.external) ?? [])],
    platform: 'node',
    write: false,
    plugins: [vanillaExtractTransformPlugin({
      identOption
    }), ...((esbuildOptions === null || esbuildOptions === void 0 ? void 0 : esbuildOptions.plugins) ?? [])],
    absWorkingDir: cwd,
    loader: esbuildOptions === null || esbuildOptions === void 0 ? void 0 : esbuildOptions.loader,
    define: esbuildOptions === null || esbuildOptions === void 0 ? void 0 : esbuildOptions.define,
    tsconfig: esbuildOptions === null || esbuildOptions === void 0 ? void 0 : esbuildOptions.tsconfig,
    conditions: esbuildOptions === null || esbuildOptions === void 0 ? void 0 : esbuildOptions.conditions
  });
  const {
    outputFiles,
    metafile
  } = result;
  if (!outputFiles || outputFiles.length !== 1) {
    throw new Error('Invalid child compilation');
  }
  return {
    source: outputFiles[0].text,
    watchFiles: Object.keys((metafile === null || metafile === void 0 ? void 0 : metafile.inputs) || {}).map(filePath => path.join(cwd, filePath))
  };
}

exports.addFileScope = addFileScope;
exports.compile = compile;
exports.cssFileFilter = cssFileFilter;
exports.deserializeCss = deserializeCss;
exports.getPackageInfo = getPackageInfo;
exports.getSourceFromVirtualCssFile = getSourceFromVirtualCssFile;
exports.hash = hash;
exports.normalizePath = normalizePath;
exports.parseFileScope = parseFileScope;
exports.processVanillaFile = processVanillaFile;
exports.serializeCss = serializeCss;
exports.serializeVanillaModule = serializeVanillaModule;
exports.stringifyFileScope = stringifyFileScope;
exports.transform = transform;
exports.transformSync = transformSync;
exports.vanillaExtractTransformPlugin = vanillaExtractTransformPlugin;
exports.virtualCssFileFilter = virtualCssFileFilter;
