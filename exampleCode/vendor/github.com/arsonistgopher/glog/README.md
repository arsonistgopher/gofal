# glog
Glog: A package to de-couple logging.

This package offers some light de-coupling to logging packages.

Dave Cheney references the tightly coupled package usage issue in his blog post: [lets-talk-about-logging](https://dave.cheney.net/2015/11/05/lets-talk-about-logging). 

Whilst this approach couples you to a logging type and set of simple methods, it does not do anything fancy and does not lock you to a package.

It offers four methods (Debug/Info/Error/Critical) and does not have any package level vars.

__Simple to Use__
1.  Import a logging package in your main code.
2.  Create a new logging var from a concrete type using your chosen logging package.
3.  Initialise the new logging var with flags and relevant setup.
4.  Instrument your code with glog.Info("message") and other methods exposed.

Done.

See example directory for use examples.

__Driver__

I wrote this because it was a pattern I started to use, so figured instead of replicating, I would build it in to a package.

My logging is now easier.

Fun fact: *Glog* isn't Go log, it's named after it's creator, David Gee (me)

Second fun fact: No global vars.
