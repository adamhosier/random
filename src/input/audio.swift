//
//  main.swift
//  Random input - audio 
//
//  Created by Adam Hosier on 2017-0407.
//  Copyright © 2017 Adam Hosier. All rights reserved.
//

import AVFoundation

class Audio : NSObject, AVAudioRecorderDelegate {
    
    var recorder : AVAudioRecorder
    var url : URL
    
    override init() {
        do {
            let dir = FileManager.default.urls(for: .documentDirectory, in: .userDomainMask)[0] as URL
            url = dir.appendingPathComponent("sound.m4a")
            let settings = [ AVFormatIDKey: Int(kAudioFormatMPEG4AAC),
                             AVSampleRateKey: 44100,
                             AVNumberOfChannelsKey: 1,
                             AVEncoderAudioQualityKey: AVAudioQuality.high.rawValue]
            recorder = try AVAudioRecorder(url: url, settings: settings)
            super.init()
            recorder.delegate = self
            recorder.prepareToRecord()
        } catch {
            fatalError("ERR 0")
        }
    }
    
    // begin capturing audio for a specified amount of time
    func capture(timeToCapture: UInt32) -> String {
        recorder.record()
        sleep(timeToCapture)
        recorder.stop()
        do {
            let content = try String(contentsOf: url, encoding: String.Encoding.ascii)
            let index = content.index(content.startIndex, offsetBy: 57500)
            return content.substring(from: index)
        } catch let e {
            fatalError(e.localizedDescription)
        }
    }
    
    func audioRecorderEncodeErrorDidOccur(_ recorder: AVAudioRecorder, error: Error?) {
        fatalError("ERR 1")
    }
    
    func audioRecorderDidFinishRecording(_ recorder: AVAudioRecorder, successfully flag: Bool) {
        print("Done")
    }
}

let a = Audio()
let data = a.capture(timeToCapture: 2)
print(data)
