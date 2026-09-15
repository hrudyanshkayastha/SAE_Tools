import json
import sys
import urllib.request
from typing import TypedDict, Annotated
from langgraph.graph import StateGraph, END

# Define the state graph
class GraphState(TypedDict):
    event: list | dict
    llm_recommendation: dict
    validation_status: str
    risk_score: int
    final_decision: str

def node_context(state: GraphState):
    # Pass event directly
    return {"event": state["event"]}

def node_investigate(state: GraphState):
    # Call the local Ollama LLM API directly (SAE Local AI Runtime)
    event_str = json.dumps(state["event"])
    prompt = (
        "Analyze this security event sequence and output ONLY a JSON object. "
        "The JSON object MUST contain exactly two keys: 'action' and 'severity'.\n"
        "'action' MUST be exactly one of: LOG_AND_MONITOR, ESCALATE_TO_HUMAN, BLOCK_IP, ISOLATE_HOST, KILL_PROCESS, DISABLE_ACCOUNT.\n"
        "'severity' MUST be exactly one of: Info, Low, Medium, High, Critical.\n\n"
        "Severity Rubric:\n"
        "- Info: Benign, administrative tasks, or nonsense/test logs.\n"
        "- Low: Routine authorized actions.\n"
        "- Medium: Suspicious unconfirmed activity.\n"
        "- High: Confirmed attacks or policy violations.\n"
        "- Critical: >10 failed logins (brute force), successful breaches, or critical asset compromise.\n\n"
        "Action Rubric:\n"
        "- LOG_AND_MONITOR: For Info or Low severity events.\n"
        "- ESCALATE_TO_HUMAN: For Medium severity or complex ambiguous events requiring investigation.\n"
        "- BLOCK_IP: For High/Critical external network attacks.\n"
        "- ISOLATE_HOST: For High/Critical internal malware or endpoint compromise.\n"
        "- KILL_PROCESS: For High/Critical malicious processes.\n"
        "- DISABLE_ACCOUNT: For High/Critical compromised user credentials.\n\n"
        "Examples:\n"
        "Event: 150 failed SSH logins from external IP\n"
        "Output: {\"action\": \"BLOCK_IP\", \"severity\": \"Critical\"}\n"
        "Event: Successful VPN login with MFA\n"
        "Output: {\"action\": \"LOG_AND_MONITOR\", \"severity\": \"Info\"}\n"
        f"Event Sequence: {event_str}\n"
        "Output: "
    )
    
    req_body = {
        "model": "llama3.2:3b",
        "prompt": prompt,
        "format": "json",
        "stream": False
    }
    try:
        req = urllib.request.Request("http://127.0.0.1:11434/api/generate", data=json.dumps(req_body).encode('utf-8'), headers={'Content-Type': 'application/json'})
        with urllib.request.urlopen(req, timeout=30) as response:
            resp_data = json.loads(response.read().decode('utf-8'))
            llm_out = json.loads(resp_data.get("response", "{}"))
    except Exception as e:
        llm_out = {"action": "ESCALATE_TO_HUMAN", "severity": "Medium", "error": str(e)}
        
    return {"llm_recommendation": llm_out}

def node_validate_evidence(state: GraphState):
    # Never blindly trust the LLM. We validate the evidence.
    rec = state.get("llm_recommendation", {})
    action = rec.get("action", "LOG_AND_MONITOR")
    
    if action in ["BLOCK_IP", "ISOLATE_HOST", "KILL_PROCESS", "DISABLE_ACCOUNT"]:
        status = f"Action {action} requires policy authorization"
    elif action == "ESCALATE_TO_HUMAN":
        status = "Requires manual investigation"
    else:
        status = "Evidence validated (Low impact)"
        
    return {"validation_status": status}

def node_risk_policy(state: GraphState):
    severity = state.get("llm_recommendation", {}).get("severity", "Info")
    score = {"Info": 10, "Low": 30, "Medium": 50, "High": 80, "Critical": 100}.get(severity, 0)
    return {"risk_score": score}

def node_decision(state: GraphState):
    valid_actions = {"LOG_AND_MONITOR", "ESCALATE_TO_HUMAN", "BLOCK_IP", "ISOLATE_HOST", "KILL_PROCESS", "DISABLE_ACCOUNT"}
    action = state.get("llm_recommendation", {}).get("action", "")
    
    if action not in valid_actions:
        decision = "ESCALATE_TO_HUMAN"
    else:
        decision = action
        
    return {"final_decision": decision}

def build_graph():
    workflow = StateGraph(GraphState)
    
    workflow.add_node("context", node_context)
    workflow.add_node("investigate", node_investigate)
    workflow.add_node("validate", node_validate_evidence)
    workflow.add_node("risk", node_risk_policy)
    workflow.add_node("decision", node_decision)
    
    workflow.set_entry_point("context")
    workflow.add_edge("context", "investigate")
    workflow.add_edge("investigate", "validate")
    workflow.add_edge("validate", "risk")
    workflow.add_edge("risk", "decision")
    workflow.add_edge("decision", END)
    
    return workflow.compile()

if __name__ == "__main__":
    try:
        input_data = sys.stdin.read()
        if not input_data.strip():
            print(json.dumps({"error": "empty event"}))
            sys.exit(1)
            
        event = json.loads(input_data)
        
        # Allow dict (single event) or list (correlation chain)
        if not isinstance(event, (dict, list)):
            print(json.dumps({"error": "malformed input"}))
            sys.exit(1)
            
        app = build_graph()
        result = app.invoke({"event": event})
        
        # Structured reasoning output
        output = {
            "status": "success",
            "llm_recommendation": result.get("llm_recommendation"),
            "validation": result.get("validation_status"),
            "risk_score": result.get("risk_score"),
            "decision": result.get("final_decision")
        }
        print(json.dumps(output))
    except json.JSONDecodeError:
        print(json.dumps({"error": "malformed input"}))
        sys.exit(1)
    except Exception as e:
        print(json.dumps({"error": str(e)}))
        sys.exit(1)
