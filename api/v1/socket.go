package v1













// TODO: Eventuell eine goroutine schreiben die vom refresh token handler gecalled wird
// TODO: und diese goroutine called dann eine andere goroutine die die spotify api
// TODO: called und werte berechnet und die dann in derselbigen goroutine immer checkt ob der client noch
// TODO: connected ist, wenn ja, wird durch einen channel zurueck an die vorherige gorutine die berechneten
// TODO: werte gesendet, wenn nein, wird der channel geclosed und die vorherige goroutine stopt den for loop
// TODO: und exited auch. Die werte werden von der ersten gorutine innerhalb des for loop dann ueber die
// TODO: socket ans frontend gesendet.




