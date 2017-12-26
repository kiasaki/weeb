function isFunction(obj: any) {
  return !!(obj && obj.constructor && obj.call && obj.apply);
}

declare interface NewObservable<T> {
  on?: (events: string, fn: Function) => T;
  off?: (events: string, fn: Function) => T;
  one?: (event: string, fn: Function) => T;
  trigger?: (event: string) => T;
}

declare interface Observable<T> {
  on: (events: string, fn: Function) => T;
  off: (events: string, fn: Function) => T;
  one: (event: string, fn: Function) => T;
  trigger: (event: string) => T;
}

export default function<T>(subject: NewObservable<T>): Observable<T> {
  var nextId = 1;
  var callbacks = {};

  subject.on = function(events, fn) {
    if (isFunction(fn)) {
      fn._id = typeof fn._id == "undefined" ? nextId++ : fn._id;

      events.replace(/\S+/g, function(name, pos) {
        if (pos === 0) fn.single = true;
        callbacks[name] = callbacks[name] || [];
        callbacks[name].push(fn);
      });
    }
    return subject;
  };

  subject.off = function(events, fn) {
    if (events == "*") {
      callbacks = {};
    } else {
      events.replace(/\S+/g, function(name) {
        if (!callbacks[name]) return;
        if (fn) {
          callbacks[name] = callbacks[name].filter(function(cb) {
            return cb._id != fn._id;
          });
        } else {
          callbacks[name] = [];
        }
      });
    }
    return subject;
  };

  subject.one = function(event, fn) {
    function on() {
      subject.off(event, on);
      fn.apply(subject, arguments);
    }
    return subject.on(event, on);
  };

  subject.trigger = function(event) {
    var args = [].slice.call(arguments, 1);
    var fns = callbacks[event] || [];

    for (var i in fns) {
      var fn = fns[i];
      if (!fn.busy) {
        fn.busy = 1;
        if (!fn.single) args = [event].concat(args);
        fn.apply(subject, args);
        fn.busy = 0;
      }
    }
    return subject;
  };

  return subject;
}
