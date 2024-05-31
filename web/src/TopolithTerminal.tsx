import 'xterm/css/xterm.css';
import React, {useEffect, useRef} from 'react';
import {Terminal} from 'xterm';
import {wasmSend} from './wasm';

export const TopolithTerminal: React.FC = () => {
    const terminalRef = useRef<HTMLDivElement | null>(null);
    const term = useRef<Terminal | null>(null);
    const inputBuffer = useRef<string>('');

    useEffect(() => {
        const initTerminal = async () => {
            if (terminalRef.current) {
                term.current = new Terminal();
                term.current.open(terminalRef.current);
                term.current.onData((char) => {
                    if (char === '\r') {
                        // Enter key pressed.
                        term.current?.write('\r\n');
                        handleCommand(inputBuffer.current)
                            .finally(() => {
                                inputBuffer.current = '';
                                writePrompt();
                            })
                    } else if (char === '\u007F') {
                        // Backspace key pressed.
                        if (inputBuffer.current.length > 0) {
                            inputBuffer.current = inputBuffer.current.slice(0, -1);
                            term.current?.write('\b \b');
                        }
                    } else {
                        // Regular character.
                        inputBuffer.current += char;
                        term.current?.write(char);
                    }
                });

                term.current.write('Topolith CLI\r\n');
                writePrompt();
            }
        };

        const handleCommand = async (command: string) => {
            const response = await wasmSend(command);
            const responseJson = JSON.stringify(response, null, ' ');
            term.current?.write(responseJson + '\r\n');
        };

        const writePrompt = () => {
            // We can use ANSI escape codes, apparently.
            const greenText = '\x1b[32m';
            const resetText = '\x1b[0m';
            term.current?.write(`${greenText}>>> ${resetText}`);
        };

        void initTerminal();

        return () => {
            term.current?.dispose();
        };
    }, []);

    return <div ref={terminalRef} style={{width: '100%', height: '100%'}}/>;
};
