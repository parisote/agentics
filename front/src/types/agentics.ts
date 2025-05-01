export type StateVariable = {
  name: string;
  type: "string" | "int" | "float" | "boolean" | "object";
};

export type FunctionHook = {
  type: "pre" | "post";
  name: string;
};

export type Tool = {
  name: string;
  description: string;
  params: {
    name: string;
    type: string;
  }[];
};

export type Node = {
  name: string;
  type: "agent";
  prompt: string;
  functions: FunctionHook[];
  tools: Tool[];
  branches?: string[];
};

export type Edge = {
  source: string;
  target: string;
};

export type AgenticsConfig = {
  entry: string;
  state: StateVariable[];
  nodes: Node[];
  edges: Edge[];
  metadata: Record<string, any>;
}; 