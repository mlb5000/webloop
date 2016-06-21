package main

import (
	"flag"
	"fmt"
	"github.com/mlb5000/webloop"
	"log"
	"net/http"
	"os"
	"time"
)

var appBind = flag.String("app-http", ":9090", "HTTP bind address for AngularJS app")
var staticBind = flag.String("static-http", ":9091", "HTTP bind address for static app")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "angular-static-seo demonstrates using WebLoop to generate a static HTML\n")
		fmt.Fprintf(os.Stderr, "website from a dynamic, single-page AngularJS application.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "\tangular-static-seo [options]\n\n")
		fmt.Fprintf(os.Stderr, "The options are:\n\n")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Example usage:\n\n")
		fmt.Fprintf(os.Stderr, "\tTo run the sample dynamic AngularJS app at http://localhost:9090 and\n")
		fmt.Fprintf(os.Stderr, "\tthe statically rendered website at http://localhost:9091:\n\n")
		fmt.Fprintf(os.Stderr, "\t    $ angular-static-seo http://example.com\n\n")
		fmt.Fprintf(os.Stderr, "\tTry browsing to both and comparing them. View the HTML source of the\n")
		fmt.Fprintf(os.Stderr, "\tstatic site to confirm that the page is, indeed, static HTML.\n")
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
	flag.Parse()

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Launching two HTTP servers serving equivalent content:\n")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "\tDynamic AngularJS app:     http://localhost%s\n", *appBind)
	fmt.Fprintf(os.Stderr, "\tStatically rendered site:  http://localhost%s\n", *staticBind)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Try browsing to both and comparing them. View the HTML source of the\n")
	fmt.Fprintf(os.Stderr, "static site to confirm that the page is, indeed, static HTML.\n")
	fmt.Fprintln(os.Stderr)

	appMux := http.NewServeMux()
	appMux.HandleFunc("/", serveApp)
	go start("app", *appBind, appMux)

	staticMux := http.NewServeMux()
	staticHandler := &webloop.StaticRenderer{
		TargetBaseURL: "http://localhost" + *appBind,
		WaitTimeout:   time.Second * 3,
		Log:           log.New(os.Stderr, "static: ", 0),
	}
	staticMux.Handle("/", staticHandler)
	start("static", *staticBind, staticMux)
}

func start(name, bind string, mux *http.ServeMux) {
	log.Printf("%s: Listening on %s", name, bind)
	err := http.ListenAndServe(bind, mux)
	if err != nil {
		log.Fatalf("%s: ListenAndServe: %s", name, err)
	}
}

func serveApp(w http.ResponseWriter, r *http.Request) {
	w.Write(page)
}

var page = []byte(`
<!doctype html>
<html ng-app="staticSEO">
<head>
  <meta charset="utf-8">
  <title>WebLoop angular-static-seo example</title>
  <script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.2.0-rc.3/angular.min.js"></script>
  <script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.2.0-rc.3/angular-route.min.js"></script>
</head>
<body>

<div ng-view></div>

<hr>

<p><a style="color: #777" href="https://sourcegraph.com/github.com/pasangsherpa/webloop/readme">WebLoop example: AngularJS static SEO</a></p>

<pre>
URL:         {{$location.url()}}

Params:      {{$route.current.params}}

User-Agent:  {{userAgent}}
</pre>

<script type=text/ng-template id="index.html">
  <h2>Angular static SEO example</h2>
  <p>
    This sample <a href="http://angularjs.org">AngularJS</a> application demonstrates how to use
    <a href="https://sourcegraph.com/github.com/pasangsherpa/webloop/readme">WebLoop</a> to
    generate a static, SEO-friendly site from a single-page AngularJS application.
  </p>
  <hr>
  <h1>Cities</h1>
  <p>Showing {{cities.length}} cities.</p>
  <ul>
    <li ng-repeat="city in cities">
      <a ng-href="/cities/{{city.id}}">{{city.name}}</a> (population: {{city.population}})
    </li>
  </ul>
</script>

<script type=text/ng-template id="detail.html">
  <p><a href="/cities">&laquo; Back to list of cities</a></p>
  <h1>{{city.name}}</h1>
  <table>
    <tr><th>Population:</th><td>{{city.population}}</td></tr>
  </table>
</script>

<script>
var allCities = [
  {id: 'shanghai', name: 'Shanghai', population: 17836133},
  {id: 'istanbul', name: 'Istanbul', population: 13854740},
  {id: 'karachi', name: 'Karachi', population: 12991000},
  {id: 'mumbai', name: 'Mumbai', population: 12478447},
  {id: 'moscow', name: 'Moscow', population: 11977988},
  {id: 'sao-paulo', name: 'Sao Paulo', population: 11821876},
  {id: 'beijing', name: 'Beijing', population: 11716000},
];

angular.module('staticSEO', ['ngRoute'])

.config(function($locationProvider, $routeProvider) {
  $locationProvider.html5Mode(true);

  $routeProvider
    .when('/cities/:city', {
      controller: 'CityCtrl',
      resolve: {
        city: function($q, $route, $timeout) {
          var cityID = $route.current.params.city;
          var deferred = $q.defer();
          // Simulate loading delay.
          $timeout(function() {
            var city = allCities.filter(function(city) {
              return city.id === cityID;
            })[0];
            if (city) deferred.resolve(city); 
            else deferred.reject('No city found with ID "' + cityID + '"');
          }, 250);
          return deferred.promise;
        },
      },
      templateUrl: 'detail.html',
    })
    .when('/cities', {
      controller: 'CitiesCtrl',
      templateUrl: 'index.html',
    })
    .otherwise({
      redirectTo: '/cities',
    });
})

.run(function($location, $rootScope, $route, $window) {
  $rootScope.userAgent = $window.navigator.userAgent;
  $rootScope.$location = $location;
  $rootScope.$route = $route;

  $rootScope.$on('$renderStaticReady', function() {
    $window.$renderStaticReady = true;
  });
  $rootScope.$on('$routeChangeBegin', function() {
    $window.$renderStaticReady = false;
  });
})

.controller('CitiesCtrl', function($scope, $timeout) {
  $timeout(function() {
    $scope.cities = allCities;
    $scope.$emit('$renderStaticReady');
  }, 250);
})

.controller('CityCtrl', function($scope, city) {
  $scope.city = city;
  $scope.$emit('$renderStaticReady');
})

;
</script>
</body>
</html>
`)
