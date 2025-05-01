"use client";

import { useState } from 'react';
import { Edge, Node } from '../types/agentics';
import { FiPlusCircle, FiTrash2 } from 'react-icons/fi';

type EdgeFormProps = {
  edges: Edge[];
  setEdges: (edges: Edge[]) => void;
  nodes: Node[];
};

export default function EdgeForm({ edges, setEdges, nodes }: EdgeFormProps) {
  const [source, setSource] = useState('');
  const [target, setTarget] = useState('');

  const addEdge = () => {
    if (!source || !target) return;
    
    // Prevent duplicate edges
    const edgeExists = edges.some(
      edge => edge.source === source && edge.target === target
    );
    
    if (edgeExists) return;
    
    setEdges([...edges, { source, target }]);
    setSource('');
    setTarget('');
  };

  const removeEdge = (index: number) => {
    setEdges(edges.filter((_, i) => i !== index));
  };

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-medium">Node Connections</h3>
      
      <div className="p-4 border rounded space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Source Node</label>
            <select
              value={source}
              onChange={(e) => setSource(e.target.value)}
              className="mt-1 w-full p-2 border rounded"
            >
              <option value="">Select source node</option>
              {nodes.map((node, index) => (
                <option key={index} value={node.name}>
                  {node.name}
                </option>
              ))}
            </select>
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700">Target Node</label>
            <select
              value={target}
              onChange={(e) => setTarget(e.target.value)}
              className="mt-1 w-full p-2 border rounded"
            >
              <option value="">Select target node</option>
              {nodes.map((node, index) => (
                <option key={index} value={node.name}>
                  {node.name}
                </option>
              ))}
            </select>
          </div>
        </div>
        
        <div className="flex justify-end">
          <button
            onClick={addEdge}
            className="px-4 py-2 text-white bg-blue-500 rounded hover:bg-blue-600 flex items-center gap-1"
          >
            <FiPlusCircle /> Add Connection
          </button>
        </div>
      </div>
      
      {edges.length > 0 && (
        <div className="mt-4">
          <h4 className="text-sm font-medium">Current Connections:</h4>
          <ul className="mt-2 space-y-2">
            {edges.map((edge, index) => (
              <li key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                <div className="flex items-center">
                  <span className="font-medium">{edge.source}</span>
                  <span className="mx-2">→</span>
                  <span className="font-medium">{edge.target}</span>
                </div>
                <button
                  onClick={() => removeEdge(index)}
                  className="text-red-500 hover:text-red-700"
                >
                  <FiTrash2 />
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
} 