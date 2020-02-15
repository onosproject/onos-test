------------------------- MODULE Trace -------------------------

INSTANCE IOUtils

Trace == IODeserialize("trace.log", TRUE)
