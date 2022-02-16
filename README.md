Suture
======

[![Build Status](https://travis-ci.org/thejerf/suture.png?branch=master)](https://travis-ci.org/thejerf/suture)

    import "github.com/thejerf/suture/v4"

Suture provides Erlang-ish supervisor trees for Go. "Supervisor trees" ->
"sutree" -> "suture" -> holds your code together when it's trying to die.

If you are reading this on pkg.go.dev, you should [visit the v4 docs](https://pkg.go.dev/github.com/thejerf/suture/v4).

It is intended to deal gracefully with the real failure cases that can
occur with supervision trees (such as burning all your CPU time endlessly
restarting dead services), while also making no unnecessary demands on the
"service" code, and providing hooks to perform adequate logging with in a
production environment.

[A blog post describing the design decisions](http://www.jerf.org/iri/post/2930)
is available.

This module is fairly fully covered
with [godoc](https://pkg.go.dev/github.com/thejerf/suture/v4)
including an example, usage, and everything else you might expect from a
README.md on GitHub. (DRY.)

v3 and before (which existed before go module support) documentation
is [also available](https://pkg.go.dev/github.com/thejerf/suture).

Special Thanks
--------------

Special thanks to the [Syncthing team](https://syncthing.net/), who have
been fantastic about working with me to push fixes upstream of them.

Major Versions
--------------

v4 is a rewrite to make Suture function
with [contexts](https://golang.org/pkg/context/). If you are using suture
for the first time, I recommend it. It also changes how logging works, to
get a single function from the user that is presented with a defined set of
structs, rather than requiring a number of closures from the consumer.

[suture v3](https://godoc.org/gopkg.in/thejerf/suture.v3) is the latest
version that does not feature contexts. It is still supported and getting
backported fixes as of now.

Code Signing
------------

Starting with the commit after ac7cf8591b, I will be signing this repository
with the ["jerf" keybase account](https://keybase.io/jerf). If you are viewing
this repository through GitHub, you should see the commits as showing as
"verified" in the commit view.

(Bear in mind that due to the nature of how git commit signing works, there
may be runs of unverified commits; what matters is that the top one is signed.)

Aspiration
----------

One of the big wins the Erlang community has with their pervasive OTP
support is that it makes it easy for them to distribute libraries that
easily fit into the OTP paradigm. It ought to someday be considered a good
idea to distribute libraries that provide some sort of supervisor tree
functionality out of the box. It is possible to provide this functionality
without explicitly depending on the Suture library.

Changelog
---------

suture uses semantic versioning and go modules.

* 4.0.2:
  * Add the ability to specify a handler for non-string panics to format
    them.
  * Fixed an issue where trying to close a currently-panicked service was
    having problems. (This may have leaked goroutines in other ways too.)
  * Merged a PR that addresses race conditions in the test suite. (These
    seem to have been isolated to the test suite and not have affected the
    core code.)
* 4.0.1:
  * Add a channel returned from ServeBackground that can be used to
    examine any error coming out of the supervisor once it is stopped.
  * Tweak up the docs to try to make it more clear suture's special
    error returns are checked via errors.Is when possible, addressing
    issue #51.
* 4.0:
  * Switched the entire API to be context based.
  * Switched how logging works to take a single closure that will be
    presented with a defined set of structs, rather than a set of closures
    for each event.
  * Consequently, "Stop" removed from the Service interface. A wrapper for
    old-style code is provided.
  * Services can now return errors. Errors will be included in the log
    message. Two special errors control restarting behavior:
      * ErrDoNotRestart indicates the service should not be restarted,
        but other services should be unaffected.
      * ErrTerminateTree indicates the parent service tree should be
        terminated. Supervisor trees can be configured to either continue
        terminating upwards, or terminate themselves but not continue
        propagating the termination upwards.
  * UnstoppedServiceReport calling semantics modified to allow correctly
    retrieving reports from entire trees. (Prior to 4.0, a report was
    only on the supervisor it was called on.)
* 3.0.4:
  * Fix a problem with adding services to a stopped supervisor.
* 3.0.3:
  * Implemented request in Issue #37, creating a new method StopWithReport
    on supervisors that reports what services failed to stop. While a bit
    tricky to use, see warning
    about
    [TOCTOU](https://en.wikipedia.org/wiki/Time-of-check_to_time-of-use)
    issues in the godoc, it can be useful at program tear-down time.
* 3.0.2:
  * Fixed issue #35 caused by the 3.0.1 change to panic when calling .Stop
    on an unServe()d supervisor. It needs to correctly notice that .Stop
    has been called, and not start up instead, which is the contract of the
    Service interface.
* 3.0.1:
  * Fixed issue #34: Calling supervisor.Stop() while something is trying
    to shut down a service could incorrectly report the  service failed to
    shut down.
  * Calling ".Stop()" on an unstarted supervisor now panics. This is
    superior to its previous behavior, which is hanging forever.
    This is justified by the fact that the Supervisor can't provide its
    guarantees about how services are started and stopped if it is not
    itself started and stopped correctly. Further pushing me in this
    direction is that it's fairly easy to use the Supervisor correctly.
* 3.0:
  * Added a default jitter of up to 50% on the restart intervals. While
    this is a backwards-compatible change from a source perspective, this
    does represent a non-trivial behavior change. It should generally be a
    good thing, but this is released as a major version as a warning.
* 2.0.4
  * Added option PassThroughPanics, to allow panics to propagate up through
    the supervisor.
* 2.0.3
  * Accepted PR #23, making the logging functions in the supervisor public.
  * Added a new Supervisor method RemoveAndWait, allowing you to make a
    best effort way to wait for a service to terminate.
  * Accepted PR #24, adding an optional IsCompletable interface that
    Services can implement that indicates they do not need to be restarted
    upon a normal return.
* 2.0.2
  * Fixed issue #21. gccgo doesn't like `case (<-c)`, with the parentheses.
    Of course the parens aren't doing anything useful anyhow. No behavior
    changes.
* 2.0.1
  * __Test code change only__. Addresses the possibility that one of the
    tests can spuriously fail if they run in a certain order.
* 2.0.0
  * Major version due to change to the signature of the logging methods:

    A race condition could occur when the Supervisor rendered the service
    name via fmt.Sprintf("%#v"), because fmt examines the entire object
    regardless of locks through reflection. 2.0.0 changes the supervisors
    to snapshot the Service's name once, when it is added, and to pass it
    to the logging methods.
  * Removal of use of sync/atomic due to possible brokenness in the Debian
    architecture.
* 1.1.2
  * TravisCI showed that the fix for 1.1.1 induced a deadlock in Go 1.4 and
    before.
  * If the supervisor is terminated before a service, the service goroutine
    could be orphaned trying the shutdown notification to the supervisor.
    This should no longer occur.
* 1.1.1
  * Per #14, the fix in 1.1.0 did not actually wait for the Supervisor
    to stop.
* 1.1.0
  * Per #12, Supervisor.stop now tries to wait for its children before
    returning. A careful reading of the original .Stop() contract
    says this is the correct behavior.
* 1.0.1
  * Fixed data race on the .state variable.
* 1.0.0
  * Initial release.
