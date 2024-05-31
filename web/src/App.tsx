import './index.css';
import React, {useLayoutEffect} from 'react';
import {TopolithTerminal} from './TopolithTerminal';
import {initializeWasm} from "./wasm.ts";

export const App: React.FC = () => {
    useLayoutEffect( () => {
        // This is all the Topolith logic that our CLI and app require to run.
        void initializeWasm();
    }, []);

    return (
        <div style={{ width: '100vw', height: '100vh' }}>
            <TopolithTerminal />
        </div>
    );
};
