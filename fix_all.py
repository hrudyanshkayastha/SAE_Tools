import re
content = open('backend/internal/engine/engine.go').read()

# Change correlate signature
content = content.replace('func (e *Engine) correlate(ctx context.Context, event models.OCSFFinding) {', 
                          'func (e *Engine) correlate(ctx context.Context, event models.OCSFFinding) *langgraph.GraphOutput {')

# Change trigger block inside correlate
old_trigger = '''	// Trigger criteria: High/Critical severity OR >= 3 correlated events
	if currentHighest >= 3 || len(ctxData.Events) >= 3 {
		e.triggerAI(ctxData)
		delete(e.correlations, corrKey)
	}
}'''
new_trigger = '''	// Trigger criteria: High/Critical severity OR >= 3 correlated events
	if currentHighest >= 3 || len(ctxData.Events) >= 3 {
		result := e.triggerAI(ctxData)
		delete(e.correlations, corrKey)
		return result
	}
	return nil
}'''
content = content.replace(old_trigger, new_trigger)

# Change triggerAI signature
old_sig = 'func (e *Engine) triggerAI(corr *CorrelationContext) {'
new_sig = 'func (e *Engine) triggerAI(corr *CorrelationContext) *langgraph.GraphOutput {'
content = content.replace(old_sig, new_sig)

# Change returns in triggerAI
content = content.replace('return\n', 'return nil\n')
content = content.replace('e.store.SaveResponseState(corr.ID, "EXECUTING")\n}', 'e.store.SaveResponseState(corr.ID, "EXECUTING")\n\treturn result\n}')

# Wait, there are multiple bare returns. I'll replace them manually below
