// swift-tools-version: 6.0
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "is-nomachine-running",
    products: [
        // Products define the executables and libraries a package produces, making them visible to other packages.
        .library(
            name: "is-nomachine-running",
            targets: ["is-nomachine-running"]),
    ],
    dependencies: [
        .package(url: "https://github.com/vapor/vapor.git", from: "4.106.7"),
    ],
    targets: [
        // Targets are the basic building blocks of a package, defining a module or a test suite.
        // Targets can depend on other targets in this package and products from dependencies.
        .target(
            name: "is-nomachine-running"),
        .testTarget(
            name: "is-nomachine-runningTests",
            dependencies: ["is-nomachine-running"]
        ),
    ]
)
