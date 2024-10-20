<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Event extends Model
{
    use HasFactory;

    protected $hidden = [
        'id',
        'created_at',
        'updated_at'
    ];

    protected $fillable = [
        'state',
        'on_home',
        'title',
        'geolink',
        'type',
        'location',
        'start_date',
        'end_date',
        'imageUrl',
        'ticket_purchase_url'
    ];

    public function translations()
    {
        return $this->hasMany(EventLanguage::class, 'event_id', 'id');
    }
}
