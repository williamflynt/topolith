/*
Declaring the:
  * `Go` (from `src/assets/wasm_exec.js` via Go)
  * `topolithSend` (from `src/assets/topolith.wasm`)
to avoid TypeScript errors.
 */
declare const Go: any;
declare function topolithSend(command: string): string;

export const initializeWasm = async () => {
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("topolith.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
    });
}

export const wasmSend = async (command: string) => {
    try {
        return JSON.parse(topolithSend(command));
    } catch (e) {
        console.error(e);
        return 'An error occurred while processing the response.';
    }
}
