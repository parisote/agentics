"use client";

import { useState } from 'react';
import { FiPlusCircle, FiTrash2 } from 'react-icons/fi';
import { FunctionHook, Node, Tool } from '../types/agentics';

type NodeFormProps = {
  nodes: Node[];
  setNodes: (nodes: Node[]) => void;
};

export default function NodeForm({ nodes, setNodes }: NodeFormProps) {
  const [nodeName, setNodeName] = useState('');
  const [prompt, setPrompt] = useState('');
  const [selectedNode, setSelectedNode] = useState<number | null>(null);

  // Function hooks state
  const [hookType, setHookType] = useState<'pre' | 'post'>('pre');
  const [hookName, setHookName] = useState('');

  // Branches state
  const [branch, setBranch] = useState('');
  const [branches, setBranches] = useState<string[]>([]);

  const addNode = () => {
    if (!nodeName.trim() || !prompt.trim()) return;

    const newNode: Node = {
      name: nodeName,
      type: 'agent',
      prompt,
      functions: [],
      tools: [],
      ...(branches.length > 0 ? { branches } : {})
    };

    setNodes([...nodes, newNode]);
    resetForm();
  };

  const updateNode = () => {
    if (selectedNode === null || !nodeName.trim() || !prompt.trim()) return;

    const updatedNodes = [...nodes];
    updatedNodes[selectedNode] = {
      ...updatedNodes[selectedNode],
      name: nodeName,
      prompt,
      ...(branches.length > 0 ? { branches } : {})
    };

    setNodes(updatedNodes);
    resetForm();
  };

  const selectNodeForEdit = (index: number) => {
    const node = nodes[index];
    setSelectedNode(index);
    setNodeName(node.name);
    setPrompt(node.prompt);
    setBranches(node.branches || []);
  };

  const removeNode = (index: number) => {
    setNodes(nodes.filter((_, i) => i !== index));
    if (selectedNode === index) {
      resetForm();
    }
  };

  const resetForm = () => {
    setNodeName('');
    setPrompt('');
    setSelectedNode(null);
    setBranches([]);
  };

  const addHook = (nodeIndex: number) => {
    if (!hookName.trim()) return;

    const newHook: FunctionHook = {
      type: hookType,
      name: hookName
    };

    const updatedNodes = [...nodes];
    updatedNodes[nodeIndex].functions = [
      ...updatedNodes[nodeIndex].functions,
      newHook
    ];

    setNodes(updatedNodes);
    setHookName('');
  };

  const removeHook = (nodeIndex: number, hookIndex: number) => {
    const updatedNodes = [...nodes];
    updatedNodes[nodeIndex].functions = updatedNodes[nodeIndex].functions.filter(
      (_, i) => i !== hookIndex
    );

    setNodes(updatedNodes);
  };

  const addBranch = () => {
    if (!branch.trim() || branches.includes(branch)) return;
    setBranches([...branches, branch]);
    setBranch('');
  };

  const removeBranch = (branchName: string) => {
    setBranches(branches.filter(b => b !== branchName));
  };

  return (
    <div className="space-y-6">
      <h3 className="text-lg font-medium">Agent Nodes</h3>
      
      <div className="space-y-4 p-4 border rounded">
        <div className="grid grid-cols-1 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Node Name</label>
            <input
              type="text"
              value={nodeName}
              onChange={(e) => setNodeName(e.target.value)}
              placeholder="e.g., detect_intent"
              className="mt-1 w-full p-2 border rounded"
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700">Prompt</label>
            <textarea
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
              placeholder="Enter the agent's prompt..."
              className="mt-1 w-full p-2 border rounded h-24"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700">Branches</label>
            <div className="flex gap-2 mt-1">
              <input
                type="text"
                value={branch}
                onChange={(e) => setBranch(e.target.value)}
                placeholder="Branch name"
                className="flex-1 p-2 border rounded"
              />
              <button
                onClick={addBranch}
                className="p-2 text-white bg-blue-500 rounded hover:bg-blue-600"
              >
                <FiPlusCircle />
              </button>
            </div>
            
            {branches.length > 0 && (
              <div className="mt-2 flex flex-wrap gap-2">
                {branches.map((b, i) => (
                  <div key={i} className="flex items-center bg-gray-100 rounded px-2 py-1">
                    <span>{b}</span>
                    <button
                      onClick={() => removeBranch(b)}
                      className="ml-2 text-red-500 hover:text-red-700"
                    >
                      <FiTrash2 size={14} />
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        <div className="flex justify-end space-x-2 mt-4">
          <button
            onClick={resetForm}
            className="px-4 py-2 border rounded hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
            onClick={selectedNode !== null ? updateNode : addNode}
            className="px-4 py-2 text-white bg-blue-500 rounded hover:bg-blue-600"
          >
            {selectedNode !== null ? 'Update' : 'Add'} Node
          </button>
        </div>
      </div>

      {nodes.length > 0 && (
        <div className="space-y-4">
          <h4 className="text-sm font-medium">Current Nodes:</h4>
          <div className="space-y-4">
            {nodes.map((node, nodeIndex) => (
              <div key={nodeIndex} className="p-4 border rounded bg-white">
                <div className="flex justify-between items-start mb-2">
                  <div>
                    <h5 className="text-base font-medium">{node.name}</h5>
                    <p className="text-sm text-gray-500 truncate">{node.prompt.slice(0, 100)}...</p>
                  </div>
                  <div className="space-x-2">
                    <button
                      onClick={() => selectNodeForEdit(nodeIndex)}
                      className="px-2 py-1 text-xs border rounded hover:bg-gray-50"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => removeNode(nodeIndex)}
                      className="px-2 py-1 text-xs text-white bg-red-500 rounded hover:bg-red-600"
                    >
                      Delete
                    </button>
                  </div>
                </div>
                
                {/* Functions section */}
                <div className="mt-4">
                  <h6 className="text-sm font-medium">Functions:</h6>
                  <div className="flex gap-2 mt-1">
                    <select
                      value={hookType}
                      onChange={(e) => setHookType(e.target.value as 'pre' | 'post')}
                      className="p-2 border rounded"
                    >
                      <option value="pre">pre</option>
                      <option value="post">post</option>
                    </select>
                    <input
                      type="text"
                      value={hookName}
                      onChange={(e) => setHookName(e.target.value)}
                      placeholder="Hook name"
                      className="flex-1 p-2 border rounded"
                    />
                    <button
                      onClick={() => addHook(nodeIndex)}
                      className="p-2 text-white bg-blue-500 rounded hover:bg-blue-600"
                    >
                      <FiPlusCircle />
                    </button>
                  </div>
                  
                  {node.functions.length > 0 ? (
                    <ul className="mt-2 space-y-1">
                      {node.functions.map((fn, fnIndex) => (
                        <li key={fnIndex} className="flex justify-between items-center p-2 bg-gray-50 rounded">
                          <span>
                            <span className="text-xs px-1.5 py-0.5 bg-gray-200 rounded mr-2">
                              {fn.type}
                            </span>
                            {fn.name}
                          </span>
                          <button
                            onClick={() => removeHook(nodeIndex, fnIndex)}
                            className="text-red-500 hover:text-red-700"
                          >
                            <FiTrash2 size={14} />
                          </button>
                        </li>
                      ))}
                    </ul>
                  ) : (
                    <p className="text-sm text-gray-500 mt-2">No functions added</p>
                  )}
                </div>

                {/* Branches display */}
                {node.branches && node.branches.length > 0 && (
                  <div className="mt-4">
                    <h6 className="text-sm font-medium">Branches:</h6>
                    <div className="flex flex-wrap gap-1 mt-1">
                      {node.branches.map((branch, i) => (
                        <span key={i} className="text-xs px-2 py-1 bg-blue-100 rounded">
                          {branch}
                        </span>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
} 