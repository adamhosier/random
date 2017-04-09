//
//  webcam.swift
//  Random input - webcam
//
//  Created by Adam Hosier on 2017-0308.
//  Copyright © 2017 Adam Hosier. All rights reserved.
//

import AVFoundation

class Webcam : NSObject, AVCaptureVideoDataOutputSampleBufferDelegate {
    
    var session : AVCaptureSession
    var data : String
    
    override init() {
        session = AVCaptureSession()
        data = ""
        super.init()
        
        // set up camera
        session.sessionPreset = AVCaptureSessionPresetLow
        let device = AVCaptureDevice.defaultDevice(withMediaType: AVMediaTypeVideo) as AVCaptureDevice
        do {
            let input = try AVCaptureDeviceInput(device: device)
            let output = AVCaptureVideoDataOutput()
            output.videoSettings = [(kCVPixelBufferPixelFormatTypeKey as NSString) :
                NSNumber(value: kCVPixelFormatType_24RGB as UInt32)]
            
            // configure camera session
            session.beginConfiguration()
            session.addInput(input)
            session.addOutput(output)
            session.commitConfiguration()
            
            // video queue
            let queue = DispatchQueue(label: "videoQueue")
            output.setSampleBufferDelegate(self, queue: queue);
        } catch let e {
            fatalError("Input error " + e.localizedDescription)
        }
    }
    
    // begin capturing frames for a specified amount of time
    func capture(timeToCapture: UInt32) -> String {
        session.startRunning()
        sleep(timeToCapture)
        session.stopRunning()
        return data
    }
    
    // when a frame is received
    func captureOutput(_ captureOutput: AVCaptureOutput!, didOutputSampleBuffer sampleBuffer: CMSampleBuffer!, from connection: AVCaptureConnection!) {
        let img = CMSampleBufferGetImageBuffer(sampleBuffer)!
        CVPixelBufferLockBaseAddress(img, CVPixelBufferLockFlags(rawValue: 0))
        
        // grab data from image, send to stdout
        let bytes = CVPixelBufferGetBaseAddressOfPlane(img, 0)!
        let count = CVPixelBufferGetDataSize(img)
        let str = String(bytesNoCopy: bytes, length:count, encoding: String.Encoding.ascii, freeWhenDone: false)!
        var output = str.replacingOccurrences(of: "\n", with: "")
        output = output.replacingOccurrences(of: "\t", with: "")
        output = output.replacingOccurrences(of: " ", with: "")
        print(output)
        
        CVPixelBufferUnlockBaseAddress(img, CVPixelBufferLockFlags(rawValue: 0))
    }
    
    // when a frame is dropped (should never be called)
    func captureOutput(_ captureOutput: AVCaptureOutput!, didDrop sampleBuffer: CMSampleBuffer!, from connection: AVCaptureConnection!) {
        fatalError("Frames have been dropped")
    }
}

let w = Webcam()
let data = w.capture(timeToCapture: 1) //capture for 1000 ms
print(data)
