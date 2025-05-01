"use client";

import { useState } from 'react';
import { StateVariable } from '../types/agentics';
import { FiPlusCircle, FiTrash2 } from 'react-icons/fi';

type StateVariableFormProps = {
  stateVariables: StateVariable[];
  setStateVariables: (variables: StateVariable[]) => void;
};

export default function StateVariableForm({ stateVariables, setStateVariables }: StateVariableFormProps) {
  const [name, setName] = useState('');
  const [type, setType] = useState<StateVariable['type']>('string');

  const addVariable = () => {
    if (!name.trim()) return;
    
    setStateVariables([...stateVariables, { name, type }]);
    setName('');
  };

  const removeVariable = (index: number) => {
    setStateVariables(stateVariables.filter((_, i) => i !== index));
  };

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-medium">State Variables</h3>
      
      <div className="flex gap-2">
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Variable name"
          className="flex-1 p-2 border rounded"
        />
        
        <select
          value={type}
          onChange={(e) => setType(e.target.value as StateVariable['type'])}
          className="p-2 border rounded"
        >
          <option value="string">string</option>
          <option value="int">int</option>
          <option value="float">float</option>
          <option value="boolean">boolean</option>
          <option value="object">object</option>
        </select>
        
        <button
          onClick={addVariable}
          className="p-2 text-white bg-blue-500 rounded hover:bg-blue-600"
        >
          <FiPlusCircle />
        </button>
      </div>
      
      {stateVariables.length > 0 && (
        <div className="mt-4">
          <h4 className="text-sm font-medium">Current Variables:</h4>
          <ul className="mt-2 space-y-2">
            {stateVariables.map((variable, index) => (
              <li key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                <div>
                  <span className="font-medium">{variable.name}</span>
                  <span className="ml-2 text-sm text-gray-500">({variable.type})</span>
                </div>
                <button
                  onClick={() => removeVariable(index)}
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