import re

content = open('backend/internal/engine/engine.go').read()

# Change correlate definition
content = content.replace('func (e *Engine) correlate(ctx context.Context, event models.OCSFFinding) {', 
                          'func (e *Engine) correlate(ctx context.Context, event models.OCSFFinding) *langgraph.GraphOutput {')

# Find the triggerAI block in correlate
block_old = '''	// Trigger criteria: High/Critical severity OR >= 3 correlated events
	if currentHighest >= 3 || len(ctxData.Events) >= 3 {
		e.TriggerAIWithReturn(ctxData)
		delete(e.correlations, corrKey)
	}'''
block_new = '''	// Trigger criteria: High/Critical severity OR >= 3 correlated events
	if currentHighest >= 3 || len(ctxData.Events) >= 3 {
		result, _ := e.TriggerAIWithReturn(ctxData)
		delete(e.correlations, corrKey)
		return result
	}
	return nil'''
content = content.replace(block_old, block_new)

# In case TriggerAIWithReturn returns an error we ignore it here because we just want the graphoutput for tests.

open('backend/internal/engine/engine.go', 'w').write(content)
