"use client";

import { useState, useEffect } from 'react';
import StateVariableForm from './StateVariableForm';
import NodeForm from './NodeForm';
import EdgeForm from './EdgeForm';
import { AgenticsConfig, Edge, Node, StateVariable } from '../types/agentics';

export default function JsonGenerator() {
  const [entrypoint, setEntrypoint] = useState('');
  const [stateVariables, setStateVariables] = useState<StateVariable[]>([]);
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);
  const [metadata, setMetadata] = useState<Record<string, any>>({});
  const [jsonOutput, setJsonOutput] = useState('');
  const [copied, setCopied] = useState(false);

  // Generate JSON whenever the data changes
  useEffect(() => {
    const config: AgenticsConfig = {
      entry: entrypoint,
      state: stateVariables,
      nodes,
      edges,
      metadata,
    };

    setJsonOutput(JSON.stringify(config, null, 2));
  }, [entrypoint, stateVariables, nodes, edges, metadata]);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(jsonOutput);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const downloadJson = () => {
    const blob = new Blob([jsonOutput], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'agentics_config.json';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <div className="max-w-6xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-8">Agentics Configuration Generator</h1>
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="space-y-8">
          {/* Entrypoint selection */}
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Entrypoint</h3>
            <select
              value={entrypoint}
              onChange={(e) => setEntrypoint(e.target.value)}
              className="w-full p-2 border rounded"
            >
              <option value="">Select entrypoint node</option>
              {nodes.map((node, index) => (
                <option key={index} value={node.name}>
                  {node.name}
                </option>
              ))}
            </select>
          </div>

          {/* State Variables Form */}
          <StateVariableForm
            stateVariables={stateVariables}
            setStateVariables={setStateVariables}
          />

          {/* Node Form */}
          <NodeForm nodes={nodes} setNodes={setNodes} />

          {/* Edge Form */}
          <EdgeForm edges={edges} setEdges={setEdges} nodes={nodes} />
        </div>

        <div className="sticky top-6 space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-medium">Generated JSON</h3>
            <div className="space-x-2">
              <button
                onClick={copyToClipboard}
                className="px-3 py-1 text-sm border rounded hover:bg-gray-50"
              >
                {copied ? 'Copied!' : 'Copy'}
              </button>
              <button
                onClick={downloadJson}
                className="px-3 py-1 text-sm text-white bg-blue-500 rounded hover:bg-blue-600"
              >
                Download
              </button>
            </div>
          </div>
          
          <div className="relative">
            <pre className="bg-gray-800 text-gray-200 p-4 rounded-md overflow-auto h-[800px]">
              {jsonOutput}
            </pre>
          </div>
        </div>
      </div>
    </div>
  );
} 